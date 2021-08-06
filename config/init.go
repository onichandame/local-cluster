package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configFilename = "local_cluster"
	envPrefix      = "lcluster"
)

func Init() error {
	if Config != nil {
		return errors.New("cannot init config twice")
	}
	Config = newConfig()
	viper.SetConfigName(configFilename)
	systemConfigDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(systemConfigDir)
	}
	viper.SetEnvPrefix(envPrefix)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Warn("no config file is found! The local cluster will start with env variables and default")
	}
	if err := initPortRange(); err != nil {
		return err
	}
	if err := initPaths(); err != nil {
		return err
	}
	return nil
}

func initPortRange() error {
	startKey := "port_range_start"
	endKey := "port_range_end"
	viper.SetDefault(startKey, 30000)
	viper.SetDefault(endKey, 40000)
	Config.PortRange.StartAt = viper.GetUint(startKey)
	Config.PortRange.EndAt = viper.GetUint(endKey)
	return nil
}

func initPaths() error {
	rootKey := "root"
	cacheKey := "cache"
	dbKey := "db"
	instancesKey := "instances"
	storageKey := "storage"
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	defaultRoot := filepath.Dir(exePath)
	viper.SetDefault(rootKey, defaultRoot)
	Config.Path.Root = viper.GetString(rootKey)
	os.Mkdir(Config.Path.Root, os.ModeDir)
	viper.SetDefault(cacheKey, filepath.Join(defaultRoot, "cache"))
	Config.Path.Cache = viper.GetString(cacheKey)
	os.Mkdir(Config.Path.Cache, os.ModeDir)
	viper.SetDefault(dbKey, filepath.Join(defaultRoot, "db"))
	Config.Path.DB = viper.GetString(dbKey)
	os.Mkdir(Config.Path.DB, os.ModeDir)
	viper.SetDefault(instancesKey, filepath.Join(defaultRoot, "instances"))
	Config.Path.Instances = viper.GetString(instancesKey)
	os.Mkdir(Config.Path.Instances, os.ModeDir)
	viper.SetDefault(storageKey, filepath.Join(defaultRoot, "storages"))
	Config.Path.Storage = viper.GetString(storageKey)
	os.Mkdir(Config.Path.Storage, os.ModeDir)
	return nil
}
