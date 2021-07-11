package instancegroup

import (
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
)

func setInstanceGroupStatus(igDef *model.InstanceGroup, status constants.InstanceGroupStatus) error {
	return db.Db.Model(igDef).Update("status", status).Error
}
