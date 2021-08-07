package instance

import (
	"errors"
	"fmt"

	appConstants "github.com/onichandame/local-cluster/constants/application"
	insConstants "github.com/onichandame/local-cluster/constants/instance"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Terminate(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	lock := getIL().getLock(instance.ID)
	lock.Lock()
	defer lock.Unlock()
	if err = terminate(instance); err != nil {
		panic(err)
	}
	return err
}

func terminate(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if instance.ID == 0 {
		panic(errors.New(fmt.Sprintf("instance %d is not created!", instance.ID)))
	}
	var ins model.Instance
	defer utils.RecoverFromError(&err)
	if err = db.Db.First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	switch ins.Status {
	case insConstants.RUNNING, insConstants.TERMINATING:
	default:
		panic(errors.New("can not terminate instances not in running state!"))
	}
	if err = db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.TERMINATING).Error; err != nil {
		panic(err)
	}
	var app model.Application
	if err = db.Db.First(&app, "name = ?", ins.ApplicationName).Error; err != nil {
		panic(err)
	}
	switch app.Type {
	case appConstants.LOCAL:
		lrm := getLRM()
		lrm.lock.Lock()
		defer lrm.lock.Unlock()
		if err = lrm.stop(ins.ID); err != nil {
			panic(err)
		}
	case appConstants.STATIC:
		ssm := getSSM()
		ssm.lock.Lock()
		defer ssm.lock.Unlock()
		if err = ssm.stop(ins.ID); err != nil {
			panic(err)
		}
	case appConstants.REMOTE:
	}
	if err = db.Db.Model(&ins).Where("status = ?", ins.Status).Update("status", insConstants.TERMINATED).Error; err != nil {
		panic(err)
	}
	return err
}
