package gateway

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(gwDef *model.Gateway) (err error) {
	utils.RecoverFromError(&err)
	if gwDef.ID != 0 {
		panic("cannot re-create gateway")
	}
	if err = db.Db.Create(gwDef).Error; err != nil {
		panic(err)
	}
	return err
}
