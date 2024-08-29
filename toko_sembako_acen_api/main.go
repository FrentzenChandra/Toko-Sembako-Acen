package main

import (
	"time"
	"toko_sembako_acen/config"
	"toko_sembako_acen/infra/database"
	"toko_sembako_acen/infra/logger"
	"toko_sembako_acen/migrations"
	"toko_sembako_acen/routers"

	"github.com/spf13/viper"
)

func main() {
	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Asia/Dhaka")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}
	masterDSN, replicaDSN := config.DbConfiguration()

	if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}
	//later separate migration
	migrations.Migrate()

	// // Gorilla Session (only used for localhost)
	// helpers.GothicSessionInit()

	// // Initialize Google OAuth2 (Only Use For localhost)
	// helpers.GothGoogle()

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

	// localhost Routing
	// routerLocalhost := routers.SetupRouteLocalhost()
	// logger.Fatalf("%v", routerLocalhost.Run(config.ServerConfigLocalHost()))

}
