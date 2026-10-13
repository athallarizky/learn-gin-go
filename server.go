package main

import (
	"learn-gin-go/config"
	"learn-gin-go/middlewares"
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

		admin := v1.Group("/admin")
		admin.Use(middlewares.AuthMiddleware(), middlewares.IsAdmin())
		{
			admin.DELETE("/article/:slug", routes.AdminDeleteArticle)
		}

		articles := v1.Group("/articles")
		{
			articles.GET("/", routes.GetHome)
			articles.GET("/:slug", routes.GetArticle)
			articles.GET("/tag/:tag", routes.GetArticleByTag)
			// example calling 2 middleware
			// articles.DELETE("/delete/:slug", middlewares.AuthMiddleware(), middlewares.IsAdmin(), routes.AdminDeleteArticle)

			protected := articles.Group("/")
			protected.Use(middlewares.AuthMiddleware())
			{
				protected.POST("/", routes.PostArticle)
				protected.PATCH("/:slug", routes.PatchArticle)
			}
		}

		auth := v1.Group("/auth")
		{

			auth.GET("/logout", routes.Logout)
			auth.GET("/:provider", routes.GetAuthProvider)
			auth.GET("/:provider/callback", routes.GetAuthProviderContext)
		}

		user := v1.Group("user")
		{
			user.GET("/articles", middlewares.AuthMiddleware(), routes.GetUserArticles)
		}
	}

	router.Run(":8088")
}
