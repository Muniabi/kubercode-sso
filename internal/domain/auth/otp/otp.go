package otp

import (
	"context"
	"github.com/google/uuid"
)

type OTP interface {
	GenerateCode() (string, error)
	VerifyCode(ctx context.Context, code string, user uuid.UUID) error
	StoreCode(ctx context.Context, code string, user uuid.UUID) error
	DeleteCode(ctx context.Context, user uuid.UUID) error
	GetCode(ctx context.Context, user uuid.UUID) (string, error)
}
