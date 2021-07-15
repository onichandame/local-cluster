package interfaces

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func createIF(svcDef interface{}, ifDef *model.ApplicationInterface) error {
	if insDef, ok := svcDef.(*model.Instance); ok {
		if insDef.ApplicationID != ifDef.ApplicationID {
			return errors.New("instance and interface definition must belong to the same application!")
		}
		insDef.Interfaces = append(insDef.Interfaces, model.ServiceInterface{ServiceID: insDef.ID, DefinitionID: ifDef.ID})
		if err := db.Db.Save(&insDef).Error; err != nil {
			return err
		}
		if err := register(&insDef.Interfaces[len(insDef.Interfaces)-1]); err != nil {
			return err
		}
	} else if igDef, ok := svcDef.(*model.InstanceGroup); ok {
		if igDef.ApplicationID != ifDef.ApplicationID {
			return errors.New("instance group and interface definition must belong to the same application!")
		}
		igDef.Interfaces = append(igDef.Interfaces, model.ServiceInterface{ServiceID: igDef.ID, DefinitionID: ifDef.ID})
		if err := db.Db.Save(&igDef).Error; err != nil {
			return err
		}
		if err := register(&igDef.Interfaces[len(igDef.Interfaces)-1]); err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("cannot create interface for type %s", utils.GetTypeName(svcDef)))
	}
	return nil
}
