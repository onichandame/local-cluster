package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func PrepareInterfaces(svcDef interface{}) error {
	if insDef, ok := svcDef.(*model.Instance); ok {
		if len(insDef.Interfaces) > 0 {
			return errors.New("interfaces already prepared!")
		}
		if err := db.Db.Preload("Application.Interfaces").First(insDef, insDef.ID).Error; err != nil {
			return err
		}
		for _, ifDef := range insDef.Application.Interfaces {
			if err := createIF(insDef, &ifDef); err != nil {
				return err
			}
		}
	} else if igDef, ok := svcDef.(*model.InstanceGroup); ok {
		if len(igDef.Interfaces) > 0 {
			return errors.New("interfaces already prepared!")
		}
		if err := db.Db.Preload("Application.Interfaces").First(igDef, igDef.ID).Error; err != nil {
			return err
		}
		for _, ifDef := range igDef.Application.Interfaces {
			if err := createIF(igDef, &ifDef); err != nil {
				return err
			}
		}
	} else {
		return errors.New("can only prepare interfaces for a valid service!")
	}
	return nil
}
