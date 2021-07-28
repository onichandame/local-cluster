package instance

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"gorm.io/gorm"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	manager := getRunnerManager()
	manager.lock.Lock()
	defer func() { manager.lock.Unlock() }()
	// audit instances in db
	var ins model.Instance
	if rows, err := db.Db.Model(&model.Instance{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			if err = db.Db.ScanRows(rows, &ins); err != nil {
				panic(err)
			}
			switch ins.Status {
			case constants.RUNNING, constants.TERMINATING:
				if manager.runners[ins.ID] == nil {
					if err = db.Db.Model(&ins).Update("status", constants.TERMINATED).Error; err != nil {
						panic(err)
					}
				}
			default:
				if manager.runners[ins.ID] != nil {
					if err = db.Db.Model(&ins).Update("status", constants.RUNNING).Error; err != nil {
						panic(err)
					}
				}
			}
		}
	}
	// audit runners
	for id, runner := range manager.runners {
		if err = db.Db.First(&ins, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				runner.cancel()
				delete(manager.runners, id)
			} else {
				panic(err)
			}
		} else {
			if ins.Status != constants.RUNNING {
				runner.cancel()
				delete(manager.runners, id)
				if err = db.Db.Model(&ins).Update("status", constants.TERMINATED).Error; err != nil {
					panic(err)
				}
			}
		}
	}
	// audit runtimes
	runtimes := listRuntimes()
	for _, runtime := range runtimes {
		removeRuntime := func() (err error) {
			defer utils.RecoverFromError(&err)
			if err = os.RemoveAll(getRuntimeDir(runtime)); err != nil {
				panic(err)
			}
			return err
		}
		var id uint
		if id, err = utils.StrToUint(runtime); err != nil {
			if err = removeRuntime(); err != nil {
				panic(err)
			}
		} else {
			if err = db.Db.First(&ins, id).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					if err = removeRuntime(); err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			} else {
				switch ins.Status {
				case constants.CRASHED, constants.TERMINATED:
					if err = removeRuntime(); err != nil {
						panic(err)
					}
				}
			}
		}
	}
	// start creating instances
	creatings := []model.Instance{}
	if err = db.Db.Where("status = ?", constants.CREATING).Find(&creatings).Error; err != nil {
		panic(err)
	}
	for _, ins = range creatings {
		if manager.runners[ins.ID] == nil {
			go Run(&ins)
		}
	}
	// restart crashed instances if required
	restartings := []model.Instance{}
	if err = db.Db.Where("status = ? AND restart_policy = ?", constants.CRASHED, constants.ALWAYS).Find(&restartings).Error; err != nil {
		panic(err)
	}
	for _, ins = range restartings {
		if manager.runners[ins.ID] == nil {
			go Run(&ins)
		}
	}
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
