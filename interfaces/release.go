package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func ReleaseIF(svcDef interface{}) error {
	svcIfs := make([]model.ServiceInterface, 0)
	if insDef, ok := svcDef.(*model.Instance); ok {
		if err := db.Db.Preload("Interfaces").Find(&insDef).Error; err != nil {
			return err
		}
		svcIfs = insDef.Interfaces
	} else if igDef, ok := svcDef.(*model.InstanceGroup); ok {
		if err := db.Db.Preload("Interfaces").Find(&igDef).Error; err != nil {
			return err
		}
		svcIfs = igDef.Interfaces
	} else {
		return errors.New("can only release interfaces for a valid service!")
	}

	for _, svcIf := range svcIfs {
		if err := deleteIF(&svcIf); err != nil {
			return err
		}
	}
	return nil
}
