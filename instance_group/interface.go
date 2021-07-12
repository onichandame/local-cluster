package instancegroup

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func createInterface(igDef *model.InstanceGroup, ifDef *model.ApplicationInterface) error {
	if igDef.ApplicationID != ifDef.ApplicationID {
		return errors.New("instance group and interface definition must point to the same application!")
	}
	var igIf model.ServiceInterface
	igIf.ServiceID = igDef.ID
	igIf.ServiceType = "instance_groups"
	igIf.DefinitionID = ifDef.ID
	if err := db.Db.Create(&igIf).Error; err != nil {
		return err
	}
	return nil
}
