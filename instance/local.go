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
	"gorm.io/gorm"
)

type LocalRunner struct {
	cmd    *exec.Cmd
	cancel context.CancelFunc
}

type LocalRunnerManager struct {
	lock    sync.Mutex
	runners map[uint]*LocalRunner
}

func (lrm *LocalRunnerManager) _run(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if _, ok := lrm.runners[instance.ID]; ok {
		panic(errors.New(fmt.Sprintf("instance %d already running! stop or audit first!", instance.ID)))
	}
	var ins model.Instance
	if err := db.Db.Preload("Template").Preload("Interfaces").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if ins.Template == nil {
		panic(errors.New(fmt.Sprintf("instance broken!")))
	}
	switch ins.Status {
	case insConstants.RESTARTING, insConstants.WAITING:
		var app model.Application
		if err := db.Db.Preload("LocalApplication.Specs", "platform = ? AND arch = ?", runtime.GOOS, runtime.GOARCH).Preload("LocalApplication.Interfaces").First(&app, "name = ?", ins.Template.ApplicationName).Error; err != nil {
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
			if ins.Template.Env != "" {
				envs := []string{}
				if err = json.Unmarshal([]byte(ins.Template.Env), &envs); err != nil {
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
	default:
		panic(errors.New(fmt.Sprintf("cannot run local if instance in status %s", ins.Status)))
	}
	return err
}

func (lrm *LocalRunnerManager) stop(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lrm.lock.Lock()
	defer lrm.lock.Unlock()
	err = lrm._stop(instance)
	return err
}

func (lrm *LocalRunnerManager) _stop(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if runner, ok := lrm.runners[instance.ID]; ok {
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

func (lrm *LocalRunnerManager) audit() (err error) {
	defer utils.RecoverFromError(&err)
	lrm.lock.Lock()
	defer lrm.lock.Unlock()
	var wg sync.WaitGroup
	for id := range lrm.runners {
		id := id
		wg.Add(1)
		go func() {
			defer utils.RecoverAndLog()
			defer wg.Done()
			var ins model.Instance
			if err := db.Db.First(&ins, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err := lrm._stop(&ins); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				switch ins.Status {
				case insConstants.CREATING, insConstants.WAITING, insConstants.RUNNING, insConstants.RESTARTING:
				default:
					if err := lrm._stop(&ins); err != nil {
						panic(err)
					}
				}
			}
		}()
	}
	wg.Wait()
	return err
}
