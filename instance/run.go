package instance

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/interfaces"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func RunInstance(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	// Only one instance can be starting at a moment
	Lock.Lock()
	defer func() { Lock.Unlock() }()
	// if running do not run again
	switch insDef.Status {
	case constants.CREATING, constants.CRASHED, constants.TERMINATED:
		// allowed states before run
	default:
		panic(errors.New(fmt.Sprintf("instance %d already running! If it is not, audit first", insDef.ID)))
	}
	// prepare runtime directory
	if err := setInstanceState(insDef, constants.CREATING); err != nil {
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
		setInstanceState(insDef, finalState)
	}()
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
	insDir := getInsDir(insDef)
	if err = application.GetSpec(&insDef.Application); err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, insDef.Application.Specs[0].Command, strings.Split(insDef.Application.Specs[0].Args, " ")...)
	RunnersMap[insDef.ID] = &Runner{cmd: cmd, cancel: cancel}
	cmd.Dir = insDir
	cmd.Env = append(cmd.Env, parseEnv(insDef)...)
	// prepare interfaces
	if err := interfaces.PrepareInterfaces(insDef); err != nil {
		return err
	}
	for _, insIf := range insDef.Interfaces {
		var ifDef model.ApplicationInterface
		if err := db.Db.First(&ifDef, insIf.DefinitionID).Error; err != nil {
			logrus.Error(insIf.DefinitionID)
			logrus.Error(err)
			return err
		}
		if ifDef.PortByEnv != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ifDef.PortByEnv, insIf.Port))
		}
		if ifDef.PortByArg != "" {
			cmd.Args = append(cmd.Args, ifDef.PortByArg, strconv.Itoa(int(insIf.Port)))
		}
	}
	if err = cmd.Start(); err != nil {
		return err
	}
	handleExit(insDef)
	return nil
}

func handleExit(insDef *model.Instance) {
	runner, ok := RunnersMap[insDef.ID]
	delete(RunnersMap, insDef.ID)
	isTerminating := insDef.Status == constants.TERMINATING
	handler := func() {
		if err := runner.cmd.Wait(); err != nil {
			logrus.Debug(err)
		}
		if isTerminating {
			logrus.Debugf("instance %d terminated", insDef.ID)
		} else {
			logrus.Warnf("instance %d crashed", insDef.ID)
			setInstanceState(insDef, constants.CRASHED)
			switch insDef.RestartPolicy {
			case constants.ALWAYS, constants.ONFAILURE:
				go RunInstance(insDef)
			}
		}
	}
	if ok {
		go handler()
	}
}

func parseEnv(insDef *model.Instance) []string {
	return strings.Split(insDef.Env, " ")
}
