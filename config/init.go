package config

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	configFilename = "local_cluster"
	envPrefix      = "lcluster"
)

type config struct {
	PortRange struct {
		StartAt uint
		EndAt   uint
	}
}

var Config config

func Init() error {
	if err := initPresets(); err != nil {
		return errors.New("failed to init root directories")
	}
	viper.SetConfigName(configFilename)
	systemConfigDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(systemConfigDir)
	}
	viper.AddConfigPath(ConfigPresets.RootDir)
	viper.SetEnvPrefix(envPrefix)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Warn("no config file is found! The local cluster will start with env variables and default")
	}
	initPortRange()
	return nil
}

func initPortRange() {
	startKey := "port_range_start"
	endKey := "port_range_end"
	viper.SetDefault(startKey, 30000)
	viper.SetDefault(endKey, 40000)
	Config.PortRange.StartAt = viper.GetUint(startKey)
	Config.PortRange.EndAt = viper.GetUint(endKey)
}
