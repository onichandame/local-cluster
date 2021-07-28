package instance

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sync"

	"github.com/onichandame/local-cluster/pkg/utils"
)

type Runner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

type RunnerManager struct {
	lock    sync.Mutex
	runners map[uint]*Runner
}

var runnerManager RunnerManager

func getRunnerManager() *RunnerManager {
	if runnerManager.runners == nil {
		runnerManager.runners = make(map[uint]*Runner)
	}
	return &runnerManager
}

func (rm *RunnerManager) run(id uint, cmd *exec.Cmd, cancel context.CancelFunc) (err error) {
	defer utils.RecoverFromError(&err)
	if rm.runners[id] != nil {
		panic(errors.New(fmt.Sprintf("instance %d already defined!", id)))
	}
	r := Runner{
		cmd:    cmd,
		cancel: cancel,
	}
	rm.runners[id] = &r
	if err = r.cmd.Start(); err != nil {
		panic(err)
	}
	return err
}
