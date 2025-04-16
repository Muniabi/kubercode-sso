package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	collection := r.db.Collection("accounts")

	// Проверяем, существует ли пользователь с таким email
	var existingUser User
	err := collection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return errors.New("user already exists")
	}

	// Устанавливаем время создания
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Создаем пользователя
	_, err = collection.InsertOne(ctx, user)
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	collection := r.db.Collection("accounts")

	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		log.Printf("[Repository] Ошибка при поиске пользователя: %v", err)
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	collection := r.db.Collection("accounts")

	var user User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *User) error {
	collection := r.db.Collection("accounts")

	user.UpdatedAt = time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *Repository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("accounts")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Token представляет собой структуру для хранения токенов
type Token struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"user_id"`
	RefreshToken string            `bson:"refresh_token"`
	ExpiresAt    time.Time         `bson:"expires_at"`
	CreatedAt    time.Time         `bson:"created_at"`
	UpdatedAt    time.Time         `bson:"updated_at"`
}

func (r *Repository) SaveToken(ctx context.Context, token *Token) error {
	collection := r.db.Collection("tokens")

	// Удаляем старые токены пользователя
	_, err := collection.DeleteMany(ctx, bson.M{"user_id": token.UserID})
	if err != nil {
		return err
	}

	// Устанавливаем время создания
	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()

	// Сохраняем новый токен
	_, err = collection.InsertOne(ctx, token)
	return err
}

func (r *Repository) GetToken(ctx context.Context, refreshToken string) (*Token, error) {
	collection := r.db.Collection("tokens")

	var token Token
	err := collection.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return &token, nil
}

func (r *Repository) DeleteToken(ctx context.Context, refreshToken string) error {
	collection := r.db.Collection("tokens")

	_, err := collection.DeleteOne(ctx, bson.M{"refresh_token": refreshToken})
	return err
}

// GetUserTokens возвращает все refresh токены пользователя
func (r *Repository) GetUserTokens(ctx context.Context, userID string) ([]*Token, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"user_id": id}
	cursor, err := r.db.Collection("tokens").Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tokens []*Token
	if err = cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

// DeleteUserTokens удаляет все токены пользователя
func (r *Repository) DeleteUserTokens(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"user_id": id}
	_, err = r.db.Collection("tokens").DeleteMany(ctx, filter)
	return err
} 