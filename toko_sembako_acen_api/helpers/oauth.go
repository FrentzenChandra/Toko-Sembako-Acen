package helpers

import (
	"toko_sembako_acen/infra/logger"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/spf13/viper"
)

func GothGoogle() {

	googleClientId := viper.GetString("OAUTH_CLIENT_ID")
	googleClientSecret := viper.GetString("OAUTH_CLIENT_SECRET")
	googleCallbackUrl := viper.GetString("OAUTH_CALLBACK_URL")


	if googleClientId == "" || googleClientSecret == "" || googleCallbackUrl == "" {
		logger.Errorf("Environment variables (CLIENT_ID, CLIENT_SECRET, CLIENT_CALLBACK_URL) are required")
	}

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, googleCallbackUrl, "email" , "profile"),
	)
}
