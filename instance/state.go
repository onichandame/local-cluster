package instance

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func initInstance(insDef *model.Instance) error {

	if err := db.Db.First(&model.Instance{}, insDef.ID).Error; err == nil {
		return errors.New(fmt.Sprintf("instance %d already created", insDef.ID))
	}
	if err := db.Db.Create(insDef).Error; err != nil {
		return err
	}
	setInstanceState(insDef, constants.CREATING)
	return nil
}

func setInstanceState(insDef *model.Instance, state constants.InstanceStatus) error {
	if err := db.Db.Where("id = ?", insDef.ID).Update("status", state).Error; err != nil {
		return err
	}
	return nil
}
