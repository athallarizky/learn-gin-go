package routes

import (
	"learn-gin-go/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

func init() {
	goth.UseProviders(
		google.New(
			config.App.GoogleClientID,
			config.App.GoogleClientSecret,
			config.App.AuthRedirectURL+"/google/callback",
			"email",
			"profile",
		),
		github.New(
			config.App.GithubClientID,
			config.App.GithubClientSecret,
			config.App.AuthRedirectURL+"/github/callback",
		))
}

func SetupAuthRoutes(router *gin.Engine) {
	// 1. Redirect user to login provider page
	router.GET("/api/v1/auth/:provider", func(c *gin.Context) {
		q := c.Request.URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = q.Encode()

		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	// 2. Callback handler after user login
	router.GET("/api/v1/auth/:provider/callback", func(c *gin.Context) {
		q := c.Request.URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   user,
		})
	})
}
