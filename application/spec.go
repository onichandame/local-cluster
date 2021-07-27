package application

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func GetSpec(appDef *model.Application) (err error) {
	defer utils.RecoverFromError(&err)
	if err := db.Db.Preload("Specs", "platform = ? AND arch = ?", runtime.GOOS, runtime.GOARCH).First(&appDef, appDef.ID).Error; err != nil {
		panic(err)
	}
	if len(appDef.Specs) < 1 {
		panic(errors.New(fmt.Sprintf("cannot find active spec for application %s", appDef.Name)))
	}
	return err
}
