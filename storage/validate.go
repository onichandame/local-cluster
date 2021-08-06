package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func validateStorage(storage *model.Storage) (err error) {
	defer func() {
		var validated bool
		if err == nil {
			validated = true
		} else {
			validated = false
		}
		if err := db.Db.Model(storage).Update("validated", validated).Error; err != nil {
			panic(err)
		}
	}()
	defer utils.RecoverFromError(&err)
	truePath := filepath.Join(config.Config.Path.Storage, strconv.Itoa(int(storage.ID)))
	if utils.PathExists(truePath) {
		if fileInfo, err := os.Stat(truePath); err != nil {
			panic(err)
		} else {
			if !fileInfo.IsDir() {
				panic(errors.New(fmt.Sprintf("storage %d cannot be initiated because the path is occupied by a file!", storage.ID)))
			}
		}
	} else {
		if err = os.MkdirAll(truePath, os.ModeDir); err != nil {
			panic(err)
		}
	}
	return err
}

func validatePath(path string) (err error) {
	storage := new(model.Storage)
	defer func() {
		var validated bool
		if err == nil {
			validated = true
		} else {
			validated = false
		}
		if storage != nil {
			if err := db.Db.Model(&storage).Update("validated", validated).Error; err != nil {
				panic(err)
			}
		}
	}()
	defer utils.RecoverFromError(&err)
	if fileInfo, err := os.Stat(path); err != nil {
		panic(err)
	} else {
		if !fileInfo.IsDir() {
			panic(errors.New(fmt.Sprintf("path %s cannot be validated as it is not a directory!", path)))
		} else {
			name := filepath.Base(path)
			if id, err := strconv.Atoi(name); err != nil {
				panic(err)
			} else {
				if err = db.Db.First(storage, uint(id)).Error; err != nil {
					panic(err)
				}
			}
		}
	}
	return err
}
