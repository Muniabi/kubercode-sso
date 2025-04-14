package mongodb

import (
	"context"
	"kubercode/internal/domain/auth"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepositoryAdapter struct {
	repo *AuthRepository
}

func NewAuthRepositoryAdapter(repo *AuthRepository) *auth.Repository {
	return &auth.Repository{
		CreateUser: func(ctx context.Context, user *auth.User) error {
			mongoUser := &User{
				ID:       user.ID,
				Email:    user.Email,
				Password: user.Password,
				IsMentor: user.IsMentor,
			}
			return repo.CreateUser(ctx, mongoUser)
		},
		GetUserByEmail: func(ctx context.Context, email string) (*auth.User, error) {
			mongoUser, err := repo.FindUserByEmail(ctx, email)
			if err != nil {
				return nil, err
			}
			if mongoUser == nil {
				return nil, nil
			}
			return &auth.User{
				ID:       mongoUser.ID,
				Email:    mongoUser.Email,
				Password: mongoUser.Password,
				IsMentor: mongoUser.IsMentor,
			}, nil
		},
		GetUserByID: func(ctx context.Context, id primitive.ObjectID) (*auth.User, error) {
			mongoUser, err := repo.FindUserByID(ctx, id)
			if err != nil {
				return nil, err
			}
			if mongoUser == nil {
				return nil, nil
			}
			return &auth.User{
				ID:       mongoUser.ID,
				Email:    mongoUser.Email,
				Password: mongoUser.Password,
				IsMentor: mongoUser.IsMentor,
			}, nil
		},
	}
} 