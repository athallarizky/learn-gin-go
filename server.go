package main

import (
	"learn-gin-go/config"
	"learn-gin-go/routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() {
	// Database Config
	config.InitDB()
	sqlDB, err := config.DB.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// Auth Config
	config.InitAuth()

	// GIN Router Config
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1 := router.Group("/api/v1/")
	{
		articles := v1.Group("/article")
		{
			articles.GET("/", routes.GetHome)
			articles.GET("/:slug", routes.GetArticle)
			articles.POST("/", routes.PostArticle)
		}

		auth := v1.Group("/auth")
		{

			auth.GET("/logout", routes.Logout)
			auth.GET("/:provider", routes.GetAuthProvider)
			auth.GET("/:provider/callback", routes.GetAuthProviderContext)
		}
	}

	router.Run(":8088")
}
