package application

import (
	"os"

	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func PrepareCache(appDef *model.Application) error {
	spec, err := GetSpec(appDef)
	if err != nil {
		return err
	}
	cachePath, err := GetCachePath(appDef)
	if err != nil {
		return err
	}
	if utils.PathExists(cachePath) {
		if spec.Hash != "" {
			if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
				os.Remove(cachePath)
				utils.Download(spec.DownloadUrl, cachePath)
				if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
					return err
				}
			}
		}
	} else {
		utils.Download(spec.DownloadUrl, cachePath)
		if spec.Hash != "" {
			if err := utils.CheckFileHash(cachePath, spec.Hash); err != nil {
				return err
			}
		}
	}
	return nil
}
