package database

import (
	"context"
	"log/slog"
	"kubercode-sso/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnector struct {
	log *slog.Logger
	cfg *config.Config
}

func NewMongoDBConnector(log *slog.Logger, cfg *config.Config) *MongoDBConnector {
	return &MongoDBConnector{log: log, cfg: cfg}
}

// Connect - подключается к базе данных и возвращает объект базы данных
func (connector *MongoDBConnector) Connect(ctx context.Context) (*mongo.Database, error) {
	connectionString := connector.cfg.MongoDBConnectionString
	
	// Настройка клиента MongoDB
	clientOptions := options.Client().
		ApplyURI(connectionString).
		SetMaxPoolSize(100).                    // Максимальный размер пула соединений
		SetMinPoolSize(5).                      // Минимальный размер пула
		SetMaxConnIdleTime(30 * time.Minute).   // Максимальное время простоя соединения
		SetConnectTimeout(10 * time.Second).    // Таймаут подключения
		SetServerSelectionTimeout(5 * time.Second) // Таймаут выбора сервера

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		connector.log.Error("Failed to connect to MongoDB", "error", err)
		return nil, err
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		connector.log.Error("Failed to ping MongoDB", "error", err)
		return nil, err
	}

	connector.log.Info("Successfully connected to MongoDB")
	
	// Получение базы данных
	db := client.Database("SSO")
	return db, nil
}

// Disconnect - функция, которая безопасно отключается от базы данных
func (connector *MongoDBConnector) Disconnect(ctx context.Context, client *mongo.Database) error {
	if err := client.Client().Disconnect(ctx); err != nil {
		connector.log.Error("Failed to disconnect from MongoDB", "error", err)
		return err
	}
	connector.log.Info("Successfully disconnected from MongoDB")
	return nil
}
