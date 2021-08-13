package instance

import (
	"errors"
	"fmt"

	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func restart(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	err = _restart(instance)
	return err
}

func _restart(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err := db.Db.First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if ins.Status != insConstants.CRASHED {
		panic(errors.New(fmt.Sprintf("cannot restart instance at status %s", ins.Status)))
	}
	if err := db.Db.Model(&ins).Update("status", insConstants.RESTARTING).Error; err != nil {
		panic(err)
	}
	if err := _run(&ins); err != nil {
		panic(err)
	}
	return err
}
