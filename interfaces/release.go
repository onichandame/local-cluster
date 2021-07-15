package interfaces

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func ReleaseIF(svcDef interface{}) error {
	Lock.Lock()
	defer func() { Lock.Unlock() }()
	if insDef, ok := svcDef.(*model.Instance); ok {
		if err := db.Db.Preload("Interfaces").Find(&insDef).Error; err != nil {
			return err
		}
		for _, insIf := range insDef.Interfaces {
			if err := deleteIF(&insIf); err != nil {
				return err
			}
		}
	} else {
		return errors.New("can only release interfaces for a valid service!")
	}
	return nil
}
