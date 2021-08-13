package instance

import (
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func fail(instance *model.Instance) (err error) {
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	err = _fail(instance)
	return err
}

func _fail(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err := db.Db.First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	if err := db.Db.Model(&ins).Update("status", insConstants.FAILED).Error; err != nil {
		panic(err)
	}
	if err := getPM().del(&ins); err != nil {
		panic(err)
	}
	return err
}
