package helpers

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
)

func GothicSessionInit() {
	key := "secret1234567890" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30      // 30 days
	isProd := false           // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store
}
