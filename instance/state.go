package instance

import (
	"errors"
	"fmt"

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
	return nil
}

func setInstanceState(insDef *model.Instance, state model.EnumValue) error {
	statuses := model.GetInstanceStatuses(db.Db)
	status, ok := statuses[state]
	if !ok {
		return errors.New(fmt.Sprintf("failed to find the status for %s", state))
	}
	if err := db.Db.Where("id = ?", insDef.ID).Update("status_id", status.ID).Error; err != nil {
		return err
	}
	return nil
}
