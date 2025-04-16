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
)

type Config struct {
	JWTSecret    string
	TokenExpiry  time.Duration
	MongoURI     string
	RedisAddr    string
	RedisPass    string
}

func main() {
	cfg := Config{
		JWTSecret:    os.Getenv("JWT_SECRET"),
		TokenExpiry:  24 * time.Hour,
		MongoURI:     os.Getenv("MONGO_URI"),
		RedisAddr:    os.Getenv("REDIS_ADDR"),
		RedisPass:    os.Getenv("REDIS_PASS"),
	}

	// Инициализация MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	// Инициализация Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       0,
	})

	// Инициализация репозитория
	authRepo := auth.NewRepository(client.Database("sso"))

	// Инициализация сервиса
	authService := auth.NewService(authRepo, cfg.JWTSecret, cfg.TokenExpiry, redisClient)

	// Инициализация обработчиков
	authHandler := handlers.NewAuthHandler(authService)

	// Инициализация роутера
	router := router.NewRouter(authHandler, authService, redisClient)

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":1488",
		Handler: router,
	}

	go func() {
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