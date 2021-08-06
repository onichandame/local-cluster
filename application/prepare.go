package application

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/onichandame/local-cluster/constants/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Prepare(appDef *model.Application) (err error) {
	defer utils.RecoverFromError(&err)
	if appDef.ID == 0 {
		panic(errors.New(("application must be created before preparing")))
	}
	switch appDef.Type {
	case application.LOCAL:
		var localApp model.LocalApplication
		if err = db.Db.Preload("Specs", "platform = ? AND arch = ?", runtime.GOOS, runtime.GOARCH).Where("application_id = ?", appDef.ID).First(&localApp).Error; err != nil {
			panic(err)
		}
		if len(localApp.Specs) < 1 {
			panic(errors.New(fmt.Sprintf("application %d doest not support the local system %s/%s", appDef.ID, runtime.GOOS, runtime.GOARCH)))
		} else if len(localApp.Specs) > 1 {
			panic(errors.New(fmt.Sprintf("application %d has conflicting definitions for the current system %s/%s. audit the specs then retry", appDef.ID, runtime.GOOS, runtime.GOARCH)))
		}
		spec := localApp.Specs[0]
		if err = Cache(appDef.ID, spec.DownloadUrl, spec.Hash); err != nil {
			panic(err)
		}
	case application.STATIC:
		var staticApp model.StaticApplication
		if err = db.Db.Where("application_id = ?", appDef.ID).First(&staticApp).Error; err != nil {
			panic(err)
		}
		if err = Cache(appDef.ID, staticApp.DownloadUrl, staticApp.Hash); err != nil {
			panic(err)
		}
	}
	return err
}
