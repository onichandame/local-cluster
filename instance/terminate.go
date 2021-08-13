package instance

import (
	"errors"

	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Terminate(instance *model.Instance) (err error) {
	err = terminate(instance)
	return err
}

func terminate(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	err = terminate(instance)
	return err
}

func _terminate(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	defer utils.RecoverFromError(&err)
	if err := db.Db.Preload("Template").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	switch ins.Status {
	case insConstants.RUNNING, insConstants.TERMINATING:
	default:
		panic(errors.New("can not terminate instances not in running state!"))
	}
	if err := db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.TERMINATING).Error; err != nil {
		panic(err)
	}
	getPM().del(&ins)
	if err := lrm.stop(&ins); err != nil {
		panic(err)
	}
	if err := ssm.stop(&ins); err != nil {
		panic(err)
	}
	if err := db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.TERMINATED).Error; err != nil {
		panic(err)
	}
	if err := db.Db.Delete(&ins).Error; err != nil {
		panic(err)
	}
	return err
}
