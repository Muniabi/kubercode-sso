package middleware

import (
	"context"
	"net/http"
	"strings"

	"kubercode/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Проверяем формат заголовка
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		resp, err := authService.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Добавляем информацию о пользователе в контекст
		c.Set("userID", resp.ID.Hex())
		c.Set("userEmail", resp.Email)
		c.Set("userIsMentor", resp.IsMentor)

		// Добавляем токен в контекст запроса
		ctx := context.WithValue(c.Request.Context(), "token", token)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
} 