package instance

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func initInstance(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if err := db.Db.First(&model.Instance{}, insDef.ID).Error; err == nil {
		panic(errors.New(fmt.Sprintf("instance %d already created", insDef.ID)))
	}
	if err := db.Db.Create(insDef).Error; err != nil {
		panic(err)
	}
	if err := setInstanceState(insDef, constants.CREATING); err != nil {
		panic(err)
	}
	return err
}

func setInstanceState(insDef *model.Instance, state constants.InstanceStatus) (err error) {
	defer utils.RecoverFromError(&err)
	if err := db.Db.Model(&model.Instance{}).Where("id = ?", insDef.ID).Update("status", state).Error; err != nil {
		panic(err)
	}
	return err
}
