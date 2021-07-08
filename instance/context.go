package instance

import (
	"context"
	"os/exec"
	"path/filepath"

	"github.com/onichandame/local-cluster/app"
	"github.com/onichandame/local-cluster/db/model"
)

var ContextMap = map[uint]context.CancelFunc{}

func runContext(insDef *model.Instance) error {
	insDir := getInsDir(insDef)
	spec, err := app.GetSpec(&insDef.Application)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	ContextMap[insDef.ID] = cancel
	cmd := exec.CommandContext(ctx, filepath.Join(insDir, spec.Entrypoint), spec.Args)
	cmd.Dir = insDir
	if err := cmd.Start(); err != nil {
		go setInstanceState(insDef, model.FAILED)
		return err
	}
	go setInstanceState(insDef, model.RUNNING)
	return nil
}

func cancelContext(insDef *model.Instance) {
	cancel, ok := ContextMap[insDef.ID]
	if ok {
		cancel()
	}
	go setInstanceState(insDef, model.FINISHED)
}
