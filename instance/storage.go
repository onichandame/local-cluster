package instance

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

func auditStorage(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	var ins model.Instance
	if err = db.Db.Preload("StorageBinding").First(&ins, instance.ID).Error; err != nil {
		panic(err)
	}
	for _, binding := range ins.StorageBindings {
		var storage model.Storage
		if err = db.Db.First(&storage, "name = ?", binding.StorageName).Error; err != nil {
			panic(err)
		}
		truePath := filepath.Join(config.Config.Path.Storage, strconv.Itoa(int(storage.ID)))
		linkPath := filepath.Join(getInsRuntimeDir(ins.ID), binding.Path)
		if pathInfo, err := os.Stat(truePath); err != nil {
			panic(err)
		} else {
			if !pathInfo.IsDir() {
				panic(errors.New(fmt.Sprintf("storage %d does not exist!", storage.ID)))
			}
		}
		link := func() {
			if err = os.Symlink(truePath, linkPath); err != nil {
				panic(err)
			}
		}
		unlink := func() {
			if err = os.Remove(linkPath); err != nil {
				panic(err)
			}
		}
		if _, err := os.Lstat(linkPath); err == nil {
			if linkedPath, err := os.Readlink(linkPath); err == nil {
				if linkedPath != truePath {
					unlink()
					link()
				}
			} else {
				panic(err)
			}
		} else if os.IsNotExist(err) {
			link()
		} else {
			panic(err)
		}
	}
	return err
}
