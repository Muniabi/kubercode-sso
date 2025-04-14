package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	db *mongo.Database
}

func NewAuthRepository(db *mongo.Database) *AuthRepository {
	return &AuthRepository{db: db}
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Email       string            `bson:"email"`
	Password    string            `bson:"password"`
	DeviceToken string            `bson:"deviceToken"`
	IsMentor    bool              `bson:"isMentor"`
	CreatedAt   time.Time         `bson:"createdAt"`
	UpdatedAt   time.Time         `bson:"updatedAt"`
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	collection := r.db.Collection("users")
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, user)
	return err
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*User, error) {
	collection := r.db.Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *User) error {
	collection := r.db.Collection("users")
	user.UpdatedAt = time.Now()

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
	return err
}

func (r *AuthRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	collection := r.db.Collection("users")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *AuthRepository) StoreRefreshToken(ctx context.Context, userID primitive.ObjectID, token string, expiresAt time.Time) error {
	collection := r.db.Collection("refresh_tokens")
	_, err := collection.InsertOne(ctx, bson.M{
		"userId":    userID,
		"token":     token,
		"expiresAt": expiresAt,
		"createdAt": time.Now(),
	})
	return err
}

func (r *AuthRepository) FindRefreshToken(ctx context.Context, token string) (*primitive.ObjectID, error) {
	collection := r.db.Collection("refresh_tokens")
	var result struct {
		UserID primitive.ObjectID `bson:"userId"`
	}
	err := collection.FindOne(ctx, bson.M{"token": token}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &result.UserID, nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	collection := r.db.Collection("refresh_tokens")
	_, err := collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (r *AuthRepository) DeleteAllUserRefreshTokens(ctx context.Context, userID primitive.ObjectID) error {
	collection := r.db.Collection("refresh_tokens")
	_, err := collection.DeleteMany(ctx, bson.M{"userId": userID})
	return err
}

func (r *AuthRepository) FindUserByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	collection := r.db.Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
} 