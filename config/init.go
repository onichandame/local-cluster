package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	configFilename = "local_config"
	envPrefix      = "lcluster"
)

func ConfigInit() {
	viper.SetConfigName(configFilename)
	systemConfigDir, err := os.UserConfigDir()
	if err == nil {
		viper.AddConfigPath(systemConfigDir)
	}
	viper.AddConfigPath(".")
	viper.SetEnvPrefix(envPrefix)
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Warn("no config file is found! The local cluster will start with env variables and default")
	}
}
