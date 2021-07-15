package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func PrepareInterfaces(svcDef interface{}) error {
	// mutex lock
	Lock.Lock()
	defer func() { Lock.Unlock() }()
	if insDef, ok := svcDef.(*model.Instance); ok {
		if len(insDef.Interfaces) > 0 {
			return errors.New("interfaces already prepared!")
		}
		ifDefs := []model.ApplicationInterface{}
		if err := db.Db.Where("application_id = ?", insDef.ApplicationID).Find(&ifDefs).Error; err != nil {
			return err
		}
		for _, ifDef := range ifDefs {
			if err := createIF(insDef, &ifDef); err != nil {
				return err
			}
		}
	} else if igDef, ok := svcDef.(*model.InstanceGroup); ok {
		if len(igDef.Interfaces) > 0 {
			return errors.New("interfaces already prepared!")
		}
		ifDefs := []model.ApplicationInterface{}
		if err := db.Db.Where("application_id = ?", igDef.ApplicationID).Find(&ifDefs).Error; err != nil {
			return err
		}
		for _, ifDef := range ifDefs {
			if err := createIF(igDef, &ifDef); err != nil {
				return err
			}
		}
	} else {
		return errors.New("can only prepare interfaces for a valid service!")
	}
	return nil
}
