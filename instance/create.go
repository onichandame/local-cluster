package instance

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	if err = db.Db.Create(instance).Error; err != nil {
		panic(err)
	}
	return err
}
