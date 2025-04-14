package repository

import (
	"context"
	"github.com/google/uuid"
	"kubercode-sso/internal/domain/auth/dto"
	"kubercode-sso/internal/domain/auth/values"
)

type AccountRepository interface {
	GetByEmail(ctx context.Context, email values.Email) (dto.UserDTO, error)
	Save(ctx context.Context, user dto.UserDTO) error
	Update(ctx context.Context, user dto.UserDTO, searchedEmail values.Email) error
	GetById(ctx context.Context, id uuid.UUID) (dto.UserDTO, error)
}
