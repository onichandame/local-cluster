package instance

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(insDef *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if insDef.ID != 0 {
		panic("cannot re-create instance")
	}
	if insDef.Service.ID == 0 && insDef.ServiceID == 0 {
		panic("service not defined!")
	}
	if insDef.ApplicationID == 0 && insDef.Application.ID == 0 {
		panic("application not defined!")
	}
	insDef.Status = constants.CREATING
	if err = db.Db.Create(insDef).Error; err != nil {
		panic(err)
	}
	return err
}
