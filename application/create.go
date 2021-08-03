package application

import (
	"errors"
	"fmt"

	"github.com/onichandame/local-cluster/constants/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(appDef *model.Application) (err error) {
	defer utils.RecoverFromError(&err)

	if appDef.ID != 0 {
		panic(errors.New(fmt.Sprintf("application %d already created", appDef.ID)))
	}
	switch appDef.Type {
	case application.LOCAL:
		if appDef.LocalApplication == nil {
			panic(errors.New("local application must be defined"))
		}
	case application.STATIC:
		if appDef.StaticApplication == nil {
			panic(errors.New("static application must be defined"))
		}
	case application.REMOTE:
		if appDef.RemoteApplication == nil {
			panic(errors.New("remote application must be defined"))
		}
	default:
		panic(errors.New(fmt.Sprintf("application type %s not recognized!", appDef.Type)))
	}
	if err = db.Db.Create(appDef).Error; err != nil {
		panic(err)
	}
	return err
}
