package instance

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

var RunnersMap = map[uint]*Runner{}

func runContext(insDef *model.Instance) error {
	insDir := getInsDir(insDef)
	spec, err := application.GetSpec(&insDef.Application)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, filepath.Join(insDir, spec.Entrypoint), spec.Args)
	RunnersMap[insDef.ID] = &Runner{cmd: cmd, cancel: cancel}
	cmd.Dir = insDir
	cmd.Env = append(cmd.Env, parseEnv(insDef)...)
	if err := cmd.Start(); err != nil {
		go setInstanceState(insDef, model.FAILED)
		return err
	}
	handleExit(insDef)
	go setInstanceState(insDef, model.RUNNING)
	return nil
}

func parseEnv(insDef *model.Instance) []string {
	return strings.Split(insDef.Env, " ")
}

func cancelContext(insDef *model.Instance) {
	runner, ok := RunnersMap[insDef.ID]
	if ok {
		delete(RunnersMap, insDef.ID)
		runner.cancel()
	}
	go setInstanceState(insDef, model.FINISHED)
}

func handleExit(insDef *model.Instance) {
	runner, ok := RunnersMap[insDef.ID]
	policies := model.GetRestartPolicies(db.Db)
	handler := func() {
		err := runner.cmd.Wait()
		if err == nil {
			setInstanceState(insDef, model.FINISHED)
			if insDef.RestartPolicyID == policies[model.ALWAYS].ID {
				go runContext(insDef)
			} else {
				delete(RunnersMap, insDef.ID)
			}
		} else {
			setInstanceState(insDef, model.FAILED)
			if utils.Contains(utils.UintSliceToIfSlice([]uint{policies[model.ALWAYS].ID, policies[model.ONFAILURE].ID}), insDef.RestartPolicyID) {
				go runContext(insDef)
			} else {
				delete(RunnersMap, insDef.ID)
			}
			logrus.Warn(err)
		}
	}
	if ok {
		go handler()
	}
}
