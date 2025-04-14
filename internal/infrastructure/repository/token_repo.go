package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"kubercode-sso/internal/domain/auth/repository"
)

// TokenRepository - интерфейс для работы с токенами
type TokenRepository struct {
	log        *slog.Logger
	collection *mongo.Collection
}

func NewTokenRepository(log *slog.Logger, db *mongo.Database) *TokenRepository {
	return &TokenRepository{
		log:        log,
		collection: db.Collection("tokens"),
	}
}

// SaveTokens - сохраняет токены в MongoDB
func (r *TokenRepository) SaveTokens(ctx context.Context, tokens repository.Token) error {
	_, err := r.collection.InsertOne(ctx, tokens)
	return err
}

// GetToken - получает токен по jti
func (r *TokenRepository) GetToken(ctx context.Context, jti uuid.UUID) (*repository.Token, error) {
	var token *repository.Token
	filter := bson.M{"_id": jti}
	result := r.collection.FindOne(ctx, filter)
	err := result.Decode(&token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *TokenRepository) GetTokens(ctx context.Context, userEmail string) ([]repository.Token, error) {
	var tokens []repository.Token
	filter := bson.M{"user_email": userEmail}

	result, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	fmt.Println(filter)
	if err = result.All(ctx, &tokens); err != nil {
		return nil, err
	}
	return tokens, nil
}

// DeleteToken - удаляет конкретный токен по его ID
func (r *TokenRepository) DeleteToken(ctx context.Context, tokenID uuid.UUID) error {
	filter := bson.M{"_id": tokenID}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

// RevokeTokens - удаляет все токены по userEmail и device_id
func (r *TokenRepository) RevokeTokens(ctx context.Context, userEmail string, deviceId string) error {
	filter := bson.M{"user_email": userEmail, "device_id": deviceId}
	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

func (r *TokenRepository) RevokeAllTokens(ctx context.Context, userEmail string) ([]repository.Token, error) {
	filter := bson.M{"user_email": userEmail}
	result, err := r.GetTokens(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	_, err = r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func (r *TokenRepository) GetAllDeviceIdFromEmail(ctx context.Context, userEmail string) ([]bson.M, error) {
	filter := bson.M{"user_email": userEmail, "token_type": "refresh"}
	projection := bson.M{"device_id": 1, "_id": 0}
	response, err := r.collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		return nil, err
	}
	var results []bson.M
	if err = response.All(ctx, &results); err != nil {
		r.log.Error("error: %v", err)
	}
	return results, nil
}
