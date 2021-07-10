package instance

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/sirupsen/logrus"
)

func RunInstance(insDef *model.Instance) error {
	// create instance data if not already
	if insDef.ID == 0 {
		if err := initInstance(insDef); err != nil {
			return err
		}
	}
	// if running do not run again
	if insDef.Status == constants.RUNNING {
		return errors.New(fmt.Sprintf("instance %d already running! If it is not, run audit instead", insDef.ID))
	}
	// prepare runtime directory
	app := model.Application{}
	if err := db.Db.First(&app, insDef.ApplicationID).Error; err != nil {
		return err
	}
	if err := application.PrepareCache(&app); err != nil {
		return err
	}
	if err := prepareRuntime(insDef); err != nil {
		return err
	}
	// prepare the cmd context
	insDir := getInsDir(insDef)
	spec, err := application.GetSpec(&app)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, spec.Entrypoint, strings.Split(spec.Args, " ")...)
	RunnersMap[insDef.ID] = &Runner{cmd: cmd, cancel: cancel}
	cmd.Dir = insDir
	cmd.Env = append(cmd.Env, parseEnv(insDef)...)
	if err := cmd.Start(); err != nil {
		go setInstanceState(insDef, constants.CRASHED)
		logrus.Warnf("failed to start instance %d", insDef.ID)
		return err
	}
	handleExit(insDef)
	go setInstanceState(insDef, constants.RUNNING)
	return nil
}

func handleExit(insDef *model.Instance) {
	runner, ok := RunnersMap[insDef.ID]
	delete(RunnersMap, insDef.ID)
	handler := func() {
		err := runner.cmd.Wait()
		if err == nil {
			setInstanceState(insDef, constants.TERMINATED)
			if insDef.RestartPolicy == constants.ALWAYS {
				go RunInstance(insDef)
			}
		} else {
			logrus.Warnf("instance %d exit with error", insDef.ID)
			logrus.Warn(err)
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
