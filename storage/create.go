package storage

import (
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Create(storage *model.Storage) (err error) {
	defer utils.RecoverFromError(&err)
	storage.Validated = false
	if err = db.Db.Create(storage).Error; err != nil {
		panic(err)
	}
	validateStorage(storage)
	return err
}
