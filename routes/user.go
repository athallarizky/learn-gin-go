package routes

import (
	"learn-gin-go/config"
	"learn-gin-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserArticles(c *gin.Context) {
	userID := c.MustGet("userID").(uint)
	var user models.User

	// item := config.DB.Where("id = ?", userID).Preload("Articles", "user_id = ?", userID).Find(&user)
	if err := config.DB.Preload("Articles").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user.Articles,
	})
}
