package routes

import (
	"errors"
	"fmt"
	"learn-gin-go/config"
	"learn-gin-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"gorm.io/gorm"
)

func GetAuthProvider(c *gin.Context) {
	q := c.Request.URL.Query()
	q.Add("provider", c.Param("provider"))
	c.Request.URL.RawQuery = q.Encode()

	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func GetAuthProviderContext(c *gin.Context) {
	provider := c.Param("provider")

	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := getOrRegisterUser(provider, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to save user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"user":   userData,
	})
}

func getOrRegisterUser(provider string, user goth.User) (*models.User, error) {
	var dbUser models.User

	result := config.DB.Where("email = ?", user.Email).First(&dbUser)

	// Handle new user
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		username := user.NickName
		if username == "" {
			username = user.Name
		}

		email := user.Email
		if email == "" {
			email = fmt.Sprintf("%s@%s.oauth", user.UserID, provider)
		}

		dbUser = models.User{
			Email:    email,
			Username: username,
			FullName: user.Name,
			SocialId: user.UserID,
			Provider: provider,
			Avatar:   user.AvatarURL,
		}

		if err := config.DB.Create(&dbUser).Error; err != nil {
			return nil, err
		}

		return &dbUser, nil
	}

	if result.Error != nil {
		return nil, result.Error
	}

	// handle return user exist
	if dbUser.Avatar == "" && user.AvatarURL != "" {
		config.DB.Model(&dbUser).Update("avatar", user.AvatarURL)
	}
	return &dbUser, nil

}

func Logout(c *gin.Context) {
	gothic.Logout(c.Writer, c.Request)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "successfully logged out",
	})
}
