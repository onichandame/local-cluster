package instance

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/onichandame/local-cluster/config"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"gorm.io/gorm"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	// audit instances in db
	if rows, err := db.Db.Model(&model.Instance{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		var wg sync.WaitGroup
		for rows.Next() {
			var ins model.Instance
			if err := db.Db.ScanRows(rows, &ins); err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer utils.RecoverAndLog()
				lock := getIL().getLock(ins.ID)
				lock.Lock()
				defer lock.Unlock()
				id := ins.ID
				var ins model.Instance
				if err := db.Db.Preload("Template").First(&ins, id).Error; err != nil {
					panic(err)
				}
				if ins.Template == nil {
					if err := _fail(&ins); err != nil {
						panic(err)
					}
					panic(errors.New(fmt.Sprintf("instance %d has no template!", ins.ID)))
				}
				if err := auditStorage(&ins); err != nil {
					panic(err)
				}
				if err := auditInsIfs(&ins); err != nil {
					panic(err)
				}
				switch ins.Status {
				case insConstants.CREATING, insConstants.RESTARTING:
					// failed to do atomic operation -> crash
					if err := _crash(&ins); err != nil {
						panic(err)
					}
				case insConstants.TERMINATING:
					// failed to terminate -> terminate
					if err := _terminate(&ins); err != nil {
						panic(err)
					}
				case insConstants.TERMINATED:
					// terminated -> (delete)
					if err := db.Db.Delete(&ins).Error; err != nil {
						panic(err)
					}
				case insConstants.CRASHED:
					// crashed -> (restart)
					if err := _restart(&ins); err != nil {
						panic(err)
					}
				case insConstants.WAITING:
					// waiting -> (crash if probe not set)
					if !getPM().has(&ins) {
						if err := _crash(&ins); err != nil {
							panic(err)
						}
					}
				}
			}()
		}
		wg.Wait()
	}
	if err = auditOrphanIfs(); err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	// audit local runners
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getLRM().audit(); err != nil {
			panic(err)
		}
	}()
	// audit static servers
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := getSSM().audit(); err != nil {
			panic(err)
		}
	}()
	// audit runtimes
	wg.Add(1)
	go func() {
		defer wg.Done()
		runtimes := listRuntimes()
		var wg sync.WaitGroup
		for _, runtime := range runtimes {
			runtime := runtime
			wg.Add(1)
			go func() {
				path := filepath.Join(config.Config.Path.Instances, runtime)
				if id, err := utils.StrToUint(runtime); err != nil {
					if err = os.RemoveAll(path); err != nil {
						panic(err)
					}
				} else {
					lock := getIL().getLock(id)
					lock.Lock()
					defer lock.Unlock()
					var ins model.Instance
					if err = db.Db.First(&ins, id).Error; err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							if err = os.RemoveAll(path); err != nil {
								panic(err)
							}
						} else {
							panic(err)
						}
					}
				}
			}()
		}
		wg.Wait()
	}()
	wg.Wait()
	return err
}

func listRuntimes() (res []string) {
	res = make([]string, 0)
	items, err := ioutil.ReadDir(config.Config.Path.Instances)
	if err != nil {
		panic(err)
	}
	for _, item := range items {
		if item.IsDir() {
			res = append(res, item.Name())
		}
	}
	return res
}
