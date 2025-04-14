package main

import (
	"context"
	"log"
	"os"
	"time"

	"kubercode/internal/domain/auth"
	"kubercode/internal/infrastructure/http/handlers"
	"kubercode/internal/infrastructure/http/middleware"
	"kubercode/internal/infrastructure/repository/mongodb"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Подключаемся к MongoDB
	ctx := context.Background()
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Проверяем подключение
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("sso")

	// Инициализируем репозиторий
	authRepo := mongodb.NewAuthRepository(db)

	// Инициализируем сервис
	authService := auth.NewService(
		authRepo,
		os.Getenv("ACCESS_SECRET"),
		os.Getenv("REFRESH_SECRET"),
		15*time.Minute,  // Access token duration
		7*24*time.Hour,  // Refresh token duration
	)

	// Инициализируем обработчики
	authHandler := handlers.NewAuthHandler(authService)

	// Создаем Gin роутер
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware())

	// API группа v1
	v1 := router.Group("/api/v1")
	{
		// Auth группа
		auth := v1.Group("/auth")
		{
			// Публичные маршруты
			auth.POST("/signup", authHandler.SignUp)
			auth.POST("/login", authHandler.Login)
			auth.POST("/verify", authHandler.VerifyToken)
			auth.POST("/restore-password", authHandler.RestorePassword)
			auth.POST("/otp/send", authHandler.SendOTP)
			auth.POST("/otp/verify", authHandler.VerifyOTP)

			// Защищенные маршруты
			protected := auth.Group("")
			protected.Use(middleware.AuthMiddleware(authService))
			{
				protected.POST("/change-password", authHandler.ChangePassword)
				protected.POST("/change-email", authHandler.ChangeEmail)
				protected.POST("/logout", authHandler.Logout)
				protected.POST("/refresh", authHandler.RefreshToken)
				protected.POST("/logout-all", authHandler.LogoutFromAllDevices)
			}
		}
	}

	// Запускаем сервер
	server := router.Run(":1488")
	if err := server; err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 