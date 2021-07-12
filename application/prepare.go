package application

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func Prepare(appDef *model.Application) error {
	if appDef.ID == 0 {
		if err := db.Db.Create(appDef).Error; err != nil {
			return err
		}
	}
	if err := PrepareCache(appDef); err != nil {
		return err
	}
	if err := WaitCache(appDef); err != nil {
		return err
	}
	return nil
}
