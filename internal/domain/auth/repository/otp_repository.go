package repository

import (
	"context"
	"github.com/google/uuid"
)

type OTPRepository interface {
	SaveOTP(ctx context.Context, code string, user uuid.UUID) error
	GetOTP(ctx context.Context, user uuid.UUID) (string, error)
	DeleteOTP(ctx context.Context, user uuid.UUID) error
}
