package otp

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"math/big"
	"kubercode-sso/config"
	"kubercode-sso/internal/domain/auth/otp"
	"kubercode-sso/internal/domain/auth/repository"
)

type OtpGenerator struct {
	otp.OTP
	cfg           *config.Config
	log           *slog.Logger
	otpRepository repository.OTPRepository
}

func NewOTPGenerator(cfg *config.Config, log *slog.Logger, otpRepository repository.OTPRepository) *OtpGenerator {
	return &OtpGenerator{
		cfg:           cfg,
		log:           log,
		otpRepository: otpRepository,
	}
}

func (o *OtpGenerator) GenerateCode() (string, error) {
	var code string
	for i := 0; i < o.cfg.OTPLength; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate OTP: %v", err)
		}
		code += n.String()
	}
	return code, nil
}
func (o *OtpGenerator) VerifyCode(ctx context.Context, code string, user uuid.UUID) error {
	result, err := o.otpRepository.GetOTP(ctx, user)
	if err != nil {
		return err
	}
	if result != code {
		return fmt.Errorf("invalid OTP code")
	}
	return nil
}

func (o *OtpGenerator) StoreCode(ctx context.Context, code string, user uuid.UUID) error {
	return o.otpRepository.SaveOTP(ctx, code, user)
}
func (o *OtpGenerator) DeleteCode(ctx context.Context, user uuid.UUID) error {
	return o.otpRepository.DeleteOTP(ctx, user)
}
func (o *OtpGenerator) GetCode(ctx context.Context, user uuid.UUID) (string, error) {
	return o.otpRepository.GetOTP(ctx, user)
}
