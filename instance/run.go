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
	"github.com/sirupsen/logrus"
)

func RunInstance(insDef *model.Instance) error {
	// create instance data if not already
	var err error
	if insDef.ID == 0 {
		if err = initInstance(insDef); err != nil {
			return err
		}
	}
	// if running do not run again
	if insDef.Status == constants.RUNNING {
		return errors.New(fmt.Sprintf("instance %d already running! If it is not, run audit instead", insDef.ID))
	}
	// prepare runtime directory
	setInstanceState(insDef, constants.CREATING)
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
	app := model.Application{}
	if err = db.Db.First(&app, insDef.ApplicationID).Error; err != nil {
		return err
	}
	if err = application.PrepareCache(&app); err != nil {
		return err
	}
	if err = prepareRuntime(insDef); err != nil {
		return err
	}
	// prepare the cmd context
	insDir := getInsDir(insDef)
	spec, err := application.GetSpec(&app)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, spec.Command, strings.Split(spec.Args, " ")...)
	RunnersMap[insDef.ID] = &Runner{cmd: cmd, cancel: cancel}
	cmd.Dir = insDir
	cmd.Env = append(cmd.Env, parseEnv(insDef)...)
	// prepare interfaces
	ifDefs := []model.ApplicationInterface{}
	if err = db.Db.Where("application_id = ?", app.ID).Find(&ifDefs).Error; err != nil {
		return err
	}
	for _, ifDef := range ifDefs {
		insIf, err := createInterface(insDef)
		if err != nil {
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
