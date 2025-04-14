package database

import (
	"github.com/redis/go-redis/v9"
	"log/slog"
	"kubercode-sso/config"
)

type RedisConnector interface {
	Connect() (*redis.Client, error)
	Close(client redis.Client) error
	ConnectToBlackListDB() (*redis.Client, error)
}

type redisConnector struct {
	log *slog.Logger
	cfg *config.Config
}

func NewRedisConnector(log *slog.Logger, cfg *config.Config) *redisConnector {
	return &redisConnector{
		log: log,
		cfg: cfg,
	}
}

func (r *redisConnector) Connect() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.cfg.RedisAddress,
		Password: r.cfg.RedisPassword,
		DB:       0,
	})
	return client, nil
}

func (r *redisConnector) ConnectToBlackListDB() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     r.cfg.RedisAddress,
		Password: r.cfg.RedisPassword,
		DB:       1,
	})
	return client, nil
}

func (r *redisConnector) Close(client redis.Client) error {
	return client.Close()
}
