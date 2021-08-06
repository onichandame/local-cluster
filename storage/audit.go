package storage

import (
	"io/ioutil"
	"path/filepath"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	// validate storages in db
	var storage model.Storage
	if rows, err := db.Db.Model(&model.Storage{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			if err = db.Db.ScanRows(rows, &storage); err != nil {
				panic(err)
			}
			validateStorage(&storage)
		}
	}

	// validate storage in fs
	if items, err := ioutil.ReadDir(config.Config.Path.Storage); err != nil {
		panic(err)
	} else {
		for _, item := range items {
			if item.IsDir() {
				validatePath(filepath.Join(config.Config.Path.Storage, item.Name()))
			}
		}
	}
	return err
}
