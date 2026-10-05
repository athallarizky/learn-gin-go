package middlewares

import (
	"net/http"
	"strings"

	"learn-gin-go/config"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  "error",
					"message": "Authorization header missing or invalid",
				},
			)
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := config.ValidateToken(tokenString)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{
					"status":  "error",
					"message": "Authorization header missing or invalid",
				},
			)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
