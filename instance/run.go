package instance

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/onichandame/local-cluster/application"
	appConstants "github.com/onichandame/local-cluster/constants/application"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func Run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if instance.ID == 0 {
		panic(errors.New("cannot run instance before create!"))
	}
	var ins model.Instance
	if err = db.Db.First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	switch ins.Status {
	case insConstants.RUNNING, insConstants.TERMINATING, insConstants.TERMINATED, insConstants.RESTARTING:
		panic(errors.New(fmt.Sprintf("cannot run instance in status %s", ins.Status)))
	}
	var app model.Application
	if err = db.Db.Preload("LocalApplication").Preload("StaticApplication").Preload("RemoteApplication").First(&app, "name = ?", instance.ApplicationName).Error; err != nil {
		panic(err)
	}
	if err = auditInsIfs(instance); err != nil {
		panic(err)
	}
	switch app.Type {
	case appConstants.LOCAL:
		lrm := getLRM()
		if err = lrm.run(&ins); err != nil {
			panic(err)
		}
	case appConstants.STATIC:
	case appConstants.REMOTE:
	}
	return err
}

func runLocal(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	// Only one instance can be starting at a moment
	manager := getRunnerManager()
	manager.lock.Lock()
	defer func() { manager.lock.Unlock() }()
	if insDef.ID == 0 {
		insDef.Status = constants.CREATING
		if err = db.Db.Create(insDef).Error; err != nil {
			panic(err)
		}
	}
	// if running do not run again
	switch insDef.Status {
	case constants.CREATING, constants.CRASHED, constants.TERMINATED:
		// allowed states before run
	default:
		panic(errors.New(fmt.Sprintf("instance %d already running! If it is not, audit first", insDef.ID)))
	}
	if manager.runners[insDef.ID] != nil {
		panic(errors.New(fmt.Sprintf("instance %d is still running!", insDef.ID)))
	}
	// init instance
	if err = setInstanceState(insDef, constants.CREATING); err != nil {
		panic(err)
	}
	defer func() {
		var finalState constants.InstanceStatus
		if err == nil {
			logrus.Infof("instance %d is now running", insDef.ID)
			finalState = constants.RUNNING
		} else {
			logrus.Warnf("instance %d failed to start", insDef.ID)
			logrus.Error(err)
			finalState = constants.CRASHED
		}
		if err = setInstanceState(insDef, finalState); err != nil {
			panic(err)
		}
	}()
	// prepare runtime
	if err = db.Db.Preload("Application").First(&insDef, insDef.ID).Error; err != nil {
		panic(err)
	}
	if err = application.PrepareCache(&insDef.Application); err != nil {
		panic(err)
	}
	if err = application.WaitCache(&insDef.Application); err != nil {
		panic(err)
	}
	if err = prepareRuntime(insDef); err != nil {
		panic(err)
	}
	// prepare the cmd context
	insDir := getInsRuntimeDir(insDef.ID)
	if err = application.GetSpec(&insDef.Application); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, insDef.Application.Specs[0].Command, strings.Split(insDef.Application.Specs[0].Args, " ")...)
	cmd.Dir = insDir
	cmd.Env = append(cmd.Env, parseEnv(insDef)...)
	// prepare interfaces
	if err = interfaces.PrepareInterfaces(insDef); err != nil {
		panic(err)
	}
	for _, insIf := range insDef.Interfaces {
		var ifDef model.ApplicationInterface
		if err = db.Db.First(&ifDef, insIf.DefinitionID).Error; err != nil {
			logrus.Error(insIf.DefinitionID)
			logrus.Error(err)
			panic(err)
		}
		if ifDef.PortByEnv != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ifDef.PortByEnv, insIf.Port))
		}
		if ifDef.PortByArg != "" {
			cmd.Args = append(cmd.Args, ifDef.PortByArg, strconv.Itoa(int(insIf.Port)))
		}
	}
	if err = cmd.Start(); err != nil {
		panic(err)
	}
	if err = manager.run(insDef.ID, cmd, cancel); err != nil {
		panic(err)
	}
	go func() {
		defer func() {
			var finalState constants.InstanceStatus
			if err := recover(); err != nil {
				finalState = constants.CRASHED
			} else {
				finalState = constants.TERMINATED
			}
			if err := setInstanceState(insDef, finalState); err != nil {
				panic(err)
			}
		}()
		runner := manager.runners[insDef.ID]
		if runner == nil {
			panic(fmt.Sprintf("failed to run instance %d", insDef.ID))
		} else {
			err := runner.cmd.Wait()
			manager.lock.Lock()
			defer func() { manager.lock.Unlock() }()
			if err := db.Db.First(insDef, insDef.ID).Error; err != nil {
				panic(err)
			}
			if err := interfaces.Release(insDef); err != nil {
				panic(err)
			}
			if insDef.Status != constants.TERMINATING {
				if err == nil {
					panic(errors.New("crashed"))
				} else {
					panic(err)
				}
			}
		}
	}()
	return err
}

func runStatic(staticIns *model.StaticApplication) (err error) {
	defer utils.RecoverFromError(&err)
	return err
}

func parseEnv(insDef *model.Instance) []string {
	return strings.Split(insDef.Env, " ")
}
