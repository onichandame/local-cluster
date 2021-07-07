package app

import (
	"os"

	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func prepareLocal(appDef *model.Application) error {
	if err := prepareCache(appDef); err != nil {
		return err
	}
	return prepareRuntime(appDef)
}

func prepareCache(appDef *model.Application) error {
	spec, err := getSpec(appDef)
	if err != nil {
		return err
	}
	cachePath, err := getCachePath(appDef)
	if err != nil {
		return err
	}
	if utils.PathExists(cachePath) {
		if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
			os.Remove(cachePath)
			utils.Download(spec.DownloadUrl, cachePath)
		}
	} else {
		utils.Download(spec.DownloadUrl, cachePath)
	}
	return nil
}

func prepareRuntime(appDef *model.Application) error {
	appDir := getAppDir(appDef)
	if !utils.PathExists(appDir) {
		return os.Mkdir(appDir, os.ModeDir)
	}
	return nil
}
