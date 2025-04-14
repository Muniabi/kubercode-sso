package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/repository"
	"time"
)

type otpRepository struct {
	repository.OTPRepository
	log         *slog.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

func NewOTPRepository(log *slog.Logger, cfg *config.Config, redisClient *redis.Client) *otpRepository {
	return &otpRepository{
		log:         log,
		cfg:         cfg,
		redisClient: redisClient,
	}
}

func (repo *otpRepository) SaveOTP(ctx context.Context, code string, user uuid.UUID) error {
	result := repo.redisClient.Set(ctx, user.String(), code, time.Second*time.Duration(repo.cfg.OTPExpirationDurationSeconds))
	repo.log.Info("result", slog.String("code", code), slog.String("user", user.String()))
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (repo *otpRepository) GetOTP(ctx context.Context, user uuid.UUID) (string, error) {
	result := repo.redisClient.Get(ctx, user.String())
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func (repo *otpRepository) DeleteOTP(ctx context.Context, user uuid.UUID) error {
	result := repo.redisClient.Del(ctx, user.String())
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
