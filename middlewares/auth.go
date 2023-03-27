package middlewares

import (
	"net/http"
	"project/auth"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Request need access token",
				"status":  http.StatusUnauthorized,
			})

			c.Abort()

			return
		}

		_, err := auth.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  err.Error(),
				"status": http.StatusUnauthorized,
			})

			c.Abort()

			return
		}

		c.Next()
	}
}

func IsLeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Request need access token",
				"status":  http.StatusUnauthorized,
			})

			c.Abort()

			return
		}

		role, err := auth.ValidateToken(tokenString)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":  err.Error(),
				"status": http.StatusUnauthorized,
			})

			c.Abort()

			return
		}

		if role != 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "Access denied. Authorization Required",
				"status":  http.StatusForbidden,
			})

			c.Abort()

			return
		}

		c.Next()
	}
}
