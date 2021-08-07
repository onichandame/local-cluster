package instance

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func prepareRuntime(instance *model.Instance) (err error) {
	defer utils.RecoverFromError(&err)
	insDir := getInsRuntimeDir(instance.ID)
	if utils.PathExists(insDir) {
		logrus.Warnf("clearing old runtime for instance %d", instance.ID)
		os.RemoveAll(insDir)
	}
	app := model.Application{}
	if err = db.Db.First(&app, "name = ?", instance.ApplicationName).Error; err != nil {
		panic(err)
	}
	cachePath := filepath.Join(config.Config.Path.Cache, strconv.Itoa(int(app.ID)))
	if utils.PathExists(insDir) {
		if err = os.RemoveAll(insDir); err != nil {
			panic(err)
		}
	}
	if !utils.PathExists(insDir) {
		if err = os.Mkdir(insDir, os.ModeDir); err != nil {
			panic(err)
		}
	}
	if err = utils.DecompressTarGZ(cachePath, insDir); err != nil {
		panic(err)
	}
	if err = auditStorage(instance); err != nil {
		panic(err)
	}
	return err
}
