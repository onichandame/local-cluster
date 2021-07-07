package app

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
)

func getCachePath(appDef *model.Application) (string, error) {
	spec, err := getSpec(appDef)
	if err != nil {
		return "", err
	}
	var cacheName string
	if url, err := http.NewRequest("GET", spec.DownloadUrl, nil); err != nil {
		return "", err
	} else {
		cacheName = path.Base(url.URL.Path)
	}
	return filepath.Join(config.ConfigPresets.CacheDir, cacheName), nil
}

func getAppDir(appDef *model.Application) string {
	return filepath.Join(config.ConfigPresets.AppsDir, appDef.Name)
}
