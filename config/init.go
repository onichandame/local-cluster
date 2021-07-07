package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	configFilename = "local_cluster"
	envPrefix      = "lcluster"
)

func ConfigInit() {
	if err := initPresets(); err != nil {
		logrus.Fatalf("failed to init root directories")
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
}
