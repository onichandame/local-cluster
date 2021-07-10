package instance

import (
	"os"

	"github.com/onichandame/local-cluster/application"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func prepareRuntime(insDef *model.Instance) error {
	insDir := getInsDir(insDef)
	if utils.PathExists(insDir) {
		logrus.Warnf("clearing old runtime for instance %d", insDef.ID)
		os.RemoveAll(insDir)
	}
	app := model.Application{}
	if err := db.Db.First(&app, insDef.ApplicationID).Error; err != nil {
		return err
	}
	cachePath := application.GetCachePath(&app)
	if !utils.PathExists(insDir) {
		if err := os.Mkdir(insDir, os.ModeDir); err != nil {
			return err
		}
	}
	if err := utils.DecompressTarGZ(cachePath, insDir); err != nil {
		return err
	}
	return nil
}
