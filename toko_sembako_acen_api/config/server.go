package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type ServerConfiguration struct {
	Port                 string
	Secret               string
	LimitCountPerRequest int64
}

func ServerConfig() string {
	viper.SetDefault("SERVER_HOST", "192.168.18.5")
	viper.SetDefault("SERVER_PORT", "9090")

	appServer := fmt.Sprintf("%s:%s", viper.GetString("SERVER_HOST"), viper.GetString("SERVER_PORT"))
	log.Print("Server Running at :", appServer)
	return appServer
}
