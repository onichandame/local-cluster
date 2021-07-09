package application

import (
	"runtime"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func GetSpec(appDef *model.Application) (*model.ApplicationSpec, error) {
	spec := model.ApplicationSpec{}
	if err := db.Db.Where("application_id = ? AND platform = ? AND arch = ?", appDef.ID, runtime.GOOS, runtime.GOARCH).First(&spec).Error; err != nil {
		return nil, err
	}
	return &spec, nil
}
