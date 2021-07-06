package config

import (
	"os"
	"path/filepath"
)

const ROOT = "lcl"

func GetRootPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

func GetCacheDir() (string, error) {
	root, err := GetRootPath()
	if err == nil {
		root = filepath.Join(root, "cache")
		os.Mkdir(root, os.ModeDir)
	}
	return root, nil
}

func GetDbDir() (string, error) {
	root, err := GetRootPath()
	if err == nil {
		root = filepath.Join(root, "db")
		os.Mkdir(root, os.ModeDir)
	}
	return root, err
}

func GetInstancesDir() (string, error) {
	root, err := GetRootPath()
	if err == nil {
		root = filepath.Join(root, "instances")
	}
	return root, err
}

func GetAppsDir() (string, error) {
	root, err := GetRootPath()
	if err == nil {
		root = filepath.Join(root, "apps")
	}
	return root, err
}
