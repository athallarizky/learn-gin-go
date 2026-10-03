package routes

import (
	"errors"
	"learn-gin-go/config"
	"learn-gin-go/models"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func GetHome(c *gin.Context) {
	items := []models.Article{}
	config.DB.Find(&items)

	c.JSON(
		200,
		gin.H{
			"status": "success",
			"data":   items,
		})
}

func GetArticle(c *gin.Context) {
	slugParam := c.Param("slug")

	var item models.Article

	result := config.DB.First(&item, "slug = ?", slugParam)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(
			404,
			gin.H{
				"status":  "error",
				"message": "record not found",
			})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   item,
	})
}

func PostArticle(c *gin.Context) {
	item := models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(
			500,
			gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		return
	}

	c.JSON(
		200,
		gin.H{
			"status":  "success",
			"message": "success create new article",
			"data":    item,
		})
}
