package config

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func InitAuth() {
	goth.UseProviders(
		google.New(
			App.GoogleClientID,
			App.GoogleClientSecret,
			App.AuthRedirectURL+"/google/callback",
			"email",
			"profile",
		),
		github.New(
			App.GithubClientID,
			App.GithubClientSecret,
			App.AuthRedirectURL+"/github/callback",
		))
}
