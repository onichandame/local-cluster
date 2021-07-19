package application

import (
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
)

func GetCachePath(appDef *model.Application) string {
	return filepath.Join(config.Config.Path.Cache, strconv.Itoa(int(appDef.ID)))
}
