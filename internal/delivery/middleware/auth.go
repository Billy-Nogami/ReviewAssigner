package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AdminToken = "admin-secret"
	UserToken  = "user-secret"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		if c.Request.Method != "GET" {
			adminToken := c.GetHeader("AdminToken")
			if adminToken != AdminToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid admin token"}})
				c.Abort()
				return
			}
		} else {
			userToken := c.GetHeader("UserToken")
			adminToken := c.GetHeader("AdminToken")
			if userToken != UserToken && adminToken != AdminToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid token"}})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}