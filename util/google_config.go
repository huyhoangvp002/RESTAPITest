package util

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetGoogleOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}
}
