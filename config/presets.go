package config

import (
	"os"
	"path/filepath"
)

type presets struct {
	RootDir      string
	CacheDir     string
	DbDir        string
	InstancesDir string
}

var ConfigPresets = presets{}

func initPresets() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	ConfigPresets.RootDir = filepath.Dir(exePath)
	ConfigPresets.CacheDir = filepath.Join(ConfigPresets.RootDir, "cache")
	os.Mkdir(ConfigPresets.CacheDir, os.ModeDir)
	ConfigPresets.DbDir = filepath.Join(ConfigPresets.RootDir, "db")
	os.Mkdir(ConfigPresets.DbDir, os.ModeDir)
	ConfigPresets.InstancesDir = filepath.Join(ConfigPresets.RootDir, "instances")
	os.Mkdir(ConfigPresets.InstancesDir, os.ModeDir)
	return nil
}
