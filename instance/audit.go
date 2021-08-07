package instance

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/chebyrash/promise"
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
		proms := []*promise.Promise{}
		for rows.Next() {
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
				var ins model.Instance
				if err = db.Db.ScanRows(rows, &ins); err != nil {
					panic(err)
				}
				lock := getIL().getLock(ins.ID)
				lock.Lock()
				defer lock.Unlock()
				if err = db.Db.First(&ins, ins.ID).Error; err != nil {
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
					if ins.Retries < ins.MaxRetries {
						if err = run(&ins); err != nil {
							panic(err)
						}
					}
				}
			}))
			return err
		}
		if _, err = promise.AllSettled(proms...).Await(); err != nil {
			panic(err)
		}
	}
	if err = auditOrphanIfs(); err != nil {
		panic(err)
	}
	proms := []*promise.Promise{}
	// audit local runners
	proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
		defer utils.SettlePromise(resolve, reject)
		lrm := getLRM()
		lrm.lock.Lock()
		defer lrm.lock.Unlock()
		proms := []*promise.Promise{}
		for id := range lrm.runners {
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
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
			}))
		}
		if _, err = promise.AllSettled(proms...).Await(); err != nil {
			panic(err)
		}
	}))
	// audit static servers
	proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
		defer utils.SettlePromise(resolve, reject)
		ssm := getSSM()
		ssm.lock.Lock()
		defer ssm.lock.Unlock()
		proms := []*promise.Promise{}
		for id := range ssm.servers {
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
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
			}))
		}
		if _, err := promise.AllSettled(proms...).Await(); err != nil {
			panic(err)
		}
	}))
	// audit runtimes
	proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
		defer utils.SettlePromise(resolve, reject)
		runtimes := listRuntimes()
		for _, runtime := range runtimes {
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
		}
	}))
	if _, err = promise.AllSettled(proms...).Await(); err != nil {
		panic(err)
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
