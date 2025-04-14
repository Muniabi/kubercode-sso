package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log/slog"
	"kubercode-sso/internal/domain/auth/dto"
	"kubercode-sso/internal/domain/auth/values"
)

type MongoAccountRepository struct {
	log        *slog.Logger
	db         *mongo.Database
	collection *mongo.Collection
}

func NewMongoAccountRepository(log *slog.Logger, db *mongo.Database) (*MongoAccountRepository, error) {
	collection := db.Collection("accounts")
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return &MongoAccountRepository{
		log:        log,
		db:         db,
		collection: collection,
	}, nil
}

// GetByEmail - находит аккаунт по email в MongoDB
func (r *MongoAccountRepository) GetByEmail(ctx context.Context, email values.Email) (dto.UserDTO, error) {
	var user dto.UserDTO
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return dto.UserDTO{}, errors.New("user not found")
	}
	return user, err
}

func (r *MongoAccountRepository) GetById(ctx context.Context, id uuid.UUID) (dto.UserDTO, error) {
	var user dto.UserDTO
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return dto.UserDTO{}, errors.New("user not found")
	}
	return user, err
}

// Save - сохраняет аккаунт в MongoDB
func (r *MongoAccountRepository) Save(ctx context.Context, user dto.UserDTO) error {
	r.log.Info("Saving user")
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
func (r *MongoAccountRepository) Update(ctx context.Context, user dto.UserDTO, searchedEmail values.Email) error {
	filter := bson.M{"email": searchedEmail}
	update := bson.M{
		"$set": bson.M{
			"email":    user.Email,
			"password": user.Password,
		},
	}

	// Логируем фильтр и данные для обновления
	fmt.Printf("Filter: %+v\n", filter)
	fmt.Printf("Update: %+v\n", update)

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("Update error: %v\n", err)
		return err
	}

	// Логируем количество обновленных документов
	fmt.Printf("Matched: %d, Modified: %d\n", result.MatchedCount, result.ModifiedCount)

	return nil
}
