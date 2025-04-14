package dto

import (
	"kubercode-sso/internal/domain/auth/values"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID          uuid.UUID      `bson:"_id"`
	Email       values.Email   `bson:"email"`
	Password    values.Password `bson:"password"`
	IsMentor    values.IsMentor `bson:"isMentor"`
	DeviceToken string         `bson:"deviceToken"`
}

func NewUserDTO(id uuid.UUID, email values.Email, password values.Password, isMentor values.IsMentor) *UserDTO {
	return &UserDTO{
		ID:          id,
		Email:       email,
		Password:    password,
		IsMentor:    isMentor,
	}
}
