package routes

import (
	"errors"
	"learn-gin-go/config"
	"learn-gin-go/models"
	"net/http"
	"strconv"
	"time"

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
	item := models.Article{}

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
	userID := c.MustGet("userID").(uint)
	slugParam := slug.Make(c.PostForm("title"))

	// Twin slug validator
	oldItem := models.Article{} // equal to: var oldItem models.Article
	// contoh penulisan lain
	// oldItem = &models.Article{}
	// oldItem = new(models.Article)

	if err := config.DB.First(&oldItem, "slug = ?", slugParam).Error; err == nil {
		slugParam = slugParam + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article{
		Title:  c.PostForm("title"),
		Desc:   c.PostForm("desc"),
		Tag:    c.PostForm("tag"),
		Slug:   slugParam,
		UserID: userID,
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

func GetArticleByTag(c *gin.Context) {
	tag := c.Param("tag")
	items := []models.Article{}

	config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&items)

	c.JSON(200,
		gin.H{
			"status": "success",
			"data":   items,
		})
}

func PatchArticle(c *gin.Context) {
	slugParam := c.Param("slug")
	loggedUserId := c.MustGet("userID").(uint) // get userID from JWT context

	itemPayload := models.Article{
		Title: c.PostForm("title"),
		Desc:  c.PostForm("desc"),
		Tag:   c.PostForm("tag"),
		Slug:  slug.Make(c.PostForm("title")),
	}

	// Get Old Data
	var oldItem models.Article
	if err := config.DB.First(&oldItem, "slug = ?", slugParam).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "error",
				"message": "Article not found",
			})
		return
	}

	if oldItem.UserID != loggedUserId {
		c.JSON(
			http.StatusForbidden,
			gin.H{
				"status":  "error",
				"message": "you are not allowed to edit this article",
			})
		return
	}

	config.DB.Model(&oldItem).Updates(itemPayload)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Article updated successfully",
		"data":    oldItem,
	})
}
