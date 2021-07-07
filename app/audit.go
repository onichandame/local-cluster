package app

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func AppAudit(appDef *model.Application) error {
	var instances []model.Instance
	if err := db.Db.Where("application_id = ?", appDef.ID).Find(&instances).Error; err != nil {
		return err
	}
	return nil
}
