package application

import (
	"os"

	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
)

func PrepareCache(appDef *model.Application) error {
	logrus.Infof("preparing cache for app %s", appDef.Name)
	spec, err := GetSpec(appDef)
	if err != nil {
		return err
	}
	cachePath := GetCachePath(appDef)
	if utils.PathExists(cachePath) {
		if spec.Hash != "" {
			if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
				os.Remove(cachePath)
				if err := utils.Download(spec.DownloadUrl, cachePath); err != nil {
					return err
				}
				if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
					return err
				}
			}
		}
	} else {
		if err := utils.Download(spec.DownloadUrl, cachePath); err != nil {
			return err
		}
		if spec.Hash != "" {
			if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
				return err
			}
		}
	}
	logrus.Infof("prepared cache for app %s", appDef.Name)
	return nil
}
