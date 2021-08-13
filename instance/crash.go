package instance

import (
	"errors"
	"fmt"

	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func crash(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	err = _crash(instance)
	return err
}

func _crash(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err := db.Db.Preload("Template").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if ins.Template == nil {
		return fail(&ins)
	}
	switch ins.Status {
	case insConstants.WAITING, insConstants.RUNNING:
		if err := db.Db.Model(&ins).Update("status", insConstants.CRASHED).Error; err != nil {
			panic(err)
		}
		if ins.Retries < ins.Template.MaxRetries {
			if err := _restart(&ins); err != nil {
				panic(err)
			}
		}
	default:
		panic(errors.New(fmt.Sprintf("cannot crash instance at status %s", ins.Status)))
	}
	return err
}
