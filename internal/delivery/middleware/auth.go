package middleware

import (
	"strings"

	"ReviewAssigner/internal/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Публичные эндпоинты — пропускаем без проверки токена
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/auth/login" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Missing or invalid Authorization header",
				},
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"error": gin.H{
					"code":    "UNAUTHORIZED",
					"message": "Invalid or expired token",
				},
			})
			c.Abort()
			return
		}

		// Сохраняем данные пользователя в контекст (на случай, если понадобится в хендлерах)
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// Защита мутирующих операций — только admin
		if c.Request.Method != "GET" && claims.Role != "admin" {
			c.JSON(403, gin.H{
				"error": gin.H{
					"code":    "FORBIDDEN",
					"message": "Admin role required for this operation",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
