package instance

import (
	"errors"
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
			if err = db.Db.ScanRows(rows, &ins); err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer utils.RecoverAndLog()
				lock := getIL().getLock(ins.ID)
				lock.Lock()
				defer lock.Unlock()
				if err = db.Db.First(&ins, ins.ID).Error; err != nil {
					panic(err)
				}
				var template model.Template
				if err = db.Db.First(&template, ins.TemplateID).Error; err != nil {
					panic(err)
				}
				if err = auditStorage(&ins); err != nil {
					panic(err)
				}
				if err = auditInsIfs(&ins); err != nil {
					panic(err)
				}
				switch ins.Status {
				case insConstants.CREATING, insConstants.RESTARTING:
					if err = db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.CRASHED).Error; err != nil {
						panic(err)
					}
				case insConstants.TERMINATING:
					if err = terminate(&ins); err != nil {
						panic(err)
					}
				case insConstants.TERMINATED:
					if err = db.Db.Delete(&ins).Error; err != nil {
						panic(err)
					}
				case insConstants.CRASHED:
					if ins.Retries < template.MaxRetries {
						if err = run(&ins); err != nil {
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
		lrm := getLRM()
		lrm.lock.Lock()
		defer lrm.lock.Unlock()
		var wg sync.WaitGroup
		for id := range lrm.runners {
			wg.Add(1)
			go func() {
				defer utils.RecoverAndLog()
				defer wg.Done()
				il := getIL().getLock(id)
				il.Lock()
				defer il.Unlock()
				var ins model.Instance
				if err = db.Db.First(&ins, id).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = lrm.stop(id); err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				} else {
					if ins.Status != insConstants.RUNNING {
						if err = lrm.stop(id); err != nil {
							panic(err)
						}
					}
				}
			}()
		}
		wg.Wait()
	}()
	// audit static servers
	wg.Add(1)
	go func() {
		defer wg.Done()
		ssm := getSSM()
		ssm.lock.Lock()
		defer ssm.lock.Unlock()
		var wg sync.WaitGroup
		for id := range ssm.servers {
			wg.Add(1)
			go func() {
				defer utils.RecoverAndLog()
				defer wg.Done()
				il := getIL().getLock(id)
				il.Lock()
				defer il.Unlock()
				var ins model.Instance
				if err := db.Db.First(&ins, id).Error; err != nil {
					panic(err)
				}
				if ins.Status != insConstants.RUNNING {
					if err := ssm.stop(id); err != nil {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()
	}()
	// audit runtimes
	wg.Add(1)
	go func() {
		defer wg.Done()
		runtimes := listRuntimes()
		var wg sync.WaitGroup
		for _, runtime := range runtimes {
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
