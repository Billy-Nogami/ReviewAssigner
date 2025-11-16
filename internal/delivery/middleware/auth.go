package middleware

import (
	"strings"
	"ReviewAssigner/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" || strings.HasPrefix(c.Request.URL.Path, "/swagger/") || c.Request.URL.Path == "/auth/login" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Missing or invalid token"}})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": gin.H{"code": "UNAUTHORIZED", "message": "Invalid token"}})
			c.Abort()
			return
		}

		// Сохраняем claims в контексте
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// Проверяем роль для POST/PUT/DELETE
		if c.Request.Method != "GET" && claims.Role != "admin" {
			c.JSON(403, gin.H{"error": gin.H{"code": "FORBIDDEN", "message": "Admin role required"}})
			c.Abort()
			return
		}

		c.Next()
	}
}
