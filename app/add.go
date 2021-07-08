package app

import "github.com/onichandame/local-cluster/db/model"

func AppAdd(appDef *model.Application) error {
	return PrepareCache(appDef)
}
