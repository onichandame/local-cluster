package instance

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func setInstanceState(insDef *model.Instance, state constants.InstanceStatus) (err error) {
	defer utils.RecoverFromError(&err)
	if err = db.Db.Model(insDef).Update("status", state).Error; err != nil {
		panic(err)
	}
	return err
}
