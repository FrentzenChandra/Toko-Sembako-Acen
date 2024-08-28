package config

import (
	"toko_sembako_acen/infra/logger"

	"github.com/spf13/viper"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfiguration
}

// SetupConfig configuration
func SetupConfig() error {
	var configuration *Configuration

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Infof("Error to reading config file, %s", err)
		return err
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		logger.Infof("error to decode, %v", err)
		return err
	}

	return nil
}
