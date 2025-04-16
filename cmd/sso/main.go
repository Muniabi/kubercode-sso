package main

import (
	"context"
	"log"
	"os"
	"time"

	"kubercode/internal/domain/auth"
	"kubercode/internal/infrastructure/http/handlers"
	"kubercode/internal/infrastructure/http/router"

	_ "kubercode/docs" // импортируем сгенерированную документацию

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @title           KuberCode SSO API
// @version         1.0
// @description     Микросервис аутентификации и авторизации KuberCode SSO.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@kubercode.com
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:1488
// @BasePath  /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
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
	authRepo := auth.NewRepository(db)

	// Инициализируем сервис
	authService := auth.NewService(
		authRepo,
		os.Getenv("JWT_SECRET"),
		60*time.Minute,  // Token expiry - 60 minutes
	)

	// Инициализируем обработчики
	authHandler := handlers.NewAuthHandler(authService)

	// Инициализация роутера
	router := router.NewRouter(authHandler, authService)

	// Запускаем сервер
	if err := router.Run(":1488"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}