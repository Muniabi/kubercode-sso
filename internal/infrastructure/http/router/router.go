package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"kubercode/internal/domain/auth"
	"kubercode/internal/infrastructure/http/handlers"
	"kubercode/internal/infrastructure/http/middleware"
)

func NewRouter(authHandler *handlers.AuthHandler, authService *auth.Service, redis *redis.Client) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Добавляем обработку OPTIONS запросов для всех маршрутов
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API группа v1
	v1 := router.Group("/api/v1")
	{
		// Auth группа
		auth := v1.Group("/auth")
		{
			// Публичные маршруты
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/login", authHandler.Login)
			auth.POST("/restore-password", authHandler.RestorePassword)
			auth.POST("/otp/send", authHandler.SendOTP)
			auth.POST("/otp/verify", authHandler.VerifyOTP)

			// Защищенные маршруты
			protected := auth.Group("")
			protected.Use(middleware.AuthMiddleware(authService))
			{
				protected.GET("/verify", authHandler.VerifyToken)
				protected.POST("/change-password", authHandler.ChangePassword)
				protected.POST("/change-email", authHandler.ChangeEmail)
				protected.POST("/logout", authHandler.Logout)
				protected.POST("/refresh", authHandler.RefreshToken)
				protected.POST("/logout-all", authHandler.LogoutFromAllDevices)
			}
		}
	}

	return router
} 