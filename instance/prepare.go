package instance

import (
	"os"

	"github.com/onichandame/local-cluster/app"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func prepareRuntime(insDef *model.Instance) error {
	insDir := getInsDir(insDef)
	cachePath, err := app.GetCachePath(&insDef.Application)
	if err != nil {
		return err
	}
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
