package config

import (
	"os"
	"path/filepath"
)

type presets struct {
	RootDir  string
	CacheDir string
	DbDir    string
	AppsDir  string
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
	ConfigPresets.AppsDir = filepath.Join(ConfigPresets.RootDir, "apps")
	os.Mkdir(ConfigPresets.AppsDir, os.ModeDir)
	return nil
}
