package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"kubercode/internal/domain/auth"
	"kubercode/internal/infrastructure/http/handlers"
	"kubercode/internal/infrastructure/http/router"

	_ "kubercode/docs" // импортируем сгенерированную документацию
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

type Config struct {
	JWTSecret    string
	TokenExpiry  time.Duration
	MongoURI     string
	RedisAddr    string
	RedisPass    string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	// Инициализация конфигурации
	cfg := Config{
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		TokenExpiry:  24 * time.Hour,
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass:    getEnv("REDIS_PASS", ""),
	}

	log.Printf("Starting server with config: MongoDB=%s, Redis=%s", cfg.MongoURI, cfg.RedisAddr)

	// Инициализация MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Проверяем подключение к MongoDB
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB")

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	// Проверяем подключение к Redis
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without Redis...", err)
		redisClient = nil
	} else {
		log.Println("Successfully connected to Redis")
	}

	// Инициализация репозитория
	authRepo := auth.NewRepository(client.Database("sso"))

	// Инициализация сервиса
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.TokenExpiry, redisClient)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)

	// Инициализация роутера
	router := router.NewRouter(authHandler, authService, redisClient)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:    ":1488",
		Handler: router,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Server starting on port 1488")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
} 