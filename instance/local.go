package instance

// local applications

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

type LocalRunner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

type LocalRunnerManager struct {
	lock    sync.Mutex
	runners map[uint]*LocalRunner
}

func (lrm *LocalRunnerManager) run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	switch instance.Status {
	case insConstants.RESTARTING, insConstants.CREATING:
	default:
		panic(errors.New(fmt.Sprintf("instance %d already running", instance.ID)))
	}
	var ins model.Instance
	if err = db.Db.Preload("Interfaces").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	var template model.Template
	if err = db.Db.First(&template, "name = ?", ins.TemplateName).Error; err != nil {
		panic(err)
	}
	var app model.Application
	if err = db.Db.Preload("LocalApplication.Specs", "platform = ? AND arch = ?", runtime.GOOS, runtime.GOARCH).Preload("LocalApplication.Interfaces").First(&app, "name = ?", template.ApplicationName).Error; err != nil {
		panic(err)
	}
	if len(app.LocalApplication.Specs) < 1 {
		panic(errors.New(fmt.Sprintf("application %d does not support the local system %s/%s", app.ID, runtime.GOOS, runtime.GOARCH)))
	} else if len(app.LocalApplication.Specs) > 1 {
		panic(errors.New(fmt.Sprintf("application %d is broken!", app.ID)))
	}
	spec := app.LocalApplication.Specs[0]
	if lrm.runners[ins.ID] == nil {
		if err = prepareRuntime(&ins); err != nil {
			panic(err)
		}
		ctx, cancel := context.WithCancel(context.Background())
		args := []string{}
		if spec.Args != "" {
			if err = json.Unmarshal([]byte(spec.Args), &args); err != nil {
				panic(err)
			}
		}
		cmd := exec.CommandContext(ctx, spec.Command, args...)
		cmd.Dir = getInsRuntimeDir(ins.ID)
		if template.Env != "" {
			envs := []string{}
			if err = json.Unmarshal([]byte(template.Env), &envs); err != nil {
				panic(err)
			}
			cmd.Env = append(cmd.Env, envs...)
		}
		for _, insIf := range ins.Interfaces {
			var ifDef *model.LocalApplicationInterface
			for _, i := range app.LocalApplication.Interfaces {
				if i.Name == insIf.DefinitionName {
					ifDef = &i
				}
			}
			if ifDef != nil {
				if ifDef.PortByArg != "" {
					cmd.Args = append(cmd.Args, ifDef.PortByArg, strconv.Itoa(int(insIf.Port)))
				}
				if ifDef.PortByEnv != "" {
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%d", ifDef.PortByArg, insIf.Port))
				}
			}
		}
		r := LocalRunner{
			cmd:    cmd,
			cancel: cancel,
		}
		lrm.runners[ins.ID] = &r
		if err = r.cmd.Start(); err != nil {
			panic(err)
		}
		go func() {
			r.cmd.Wait()
			lrm.lock.Lock()
			defer lrm.lock.Unlock()
			delete(lrm.runners, ins.ID)
		}()
	} else {
		panic(errors.New(fmt.Sprintf("instance %d already running!", ins.ID)))
	}
	return err
}

func (lrm *LocalRunnerManager) stop(id uint) (err error) {
	defer utils.RecoverFromError(&err)
	if runner, ok := lrm.runners[id]; ok {
		runner.cancel()
	}
	return err
}

var lrm *LocalRunnerManager

func getLRM() *LocalRunnerManager {
	if lrm == nil {
		atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&lrm)), unsafe.Pointer(nil), unsafe.Pointer(&LocalRunnerManager{runners: make(map[uint]*LocalRunner)}))
	}
	return lrm
}
