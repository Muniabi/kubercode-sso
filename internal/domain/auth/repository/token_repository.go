package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type Token struct {
	Jti       uuid.UUID `bson:"_id"`
	Token     string    `bson:"token"`
	TokenType string    `bson:"token_type"`
	UserEmail string    `bson:"user_email"`
	DeviceId  string    `bson:"device_id"`
}

func NewToken(jti uuid.UUID, token string, tokenType string, userEmail string, deviceId string) *Token {
	return &Token{Jti: jti, Token: token, TokenType: tokenType, UserEmail: userEmail, DeviceId: deviceId}
}

type AbstractTokenRepository interface {
	SaveTokens(ctx context.Context, tokens Token) error
	GetTokens(ctx context.Context, userEmail string) ([]Token, error)
	RevokeTokens(ctx context.Context, userEmail string, deviceId string) error
	DeleteToken(ctx context.Context, tokenID uuid.UUID) error
	GetToken(ctx context.Context, jti uuid.UUID) (*Token, error)
	RevokeAllTokens(ctx context.Context, userEmail string) ([]Token, error)
	GetAllDeviceIdFromEmail(ctx context.Context, userEmail string) ([]bson.M, error)
}
