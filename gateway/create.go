package gateway

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(gwDef *model.Gateway) (err error) {
	utils.RecoverFromError(&err)
	if gwDef.ID != 0 {
		panic(errors.New("cannot re-create gateway"))
	}
	if gwDef.Port == 0 {
		panic(errors.New("gateway must have a definded port!"))
	}

	if err = db.Db.Create(gwDef).Error; err != nil {
		panic(err)
	}

	// find backends and start proxy
	interfaces :=[]model.ServiceInterface{}
	if err=db.Db.Where("definition_id").Find(&interfaces)
	return err
}
