package dto

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"kubercode-sso/internal/domain/auth/values"

	"github.com/google/uuid"
)

func TestNewUserDTO(t *testing.T) {
	id := uuid.New()
	email, err := values.NewEmail("test@example.com")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	password, err := values.NewPassword("password123")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	isMentor, err := values.NewIsMentor(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	userDTO := NewUserDTO(id, *email, *password, *isMentor)

	if userDTO.ID != id {
		t.Errorf("expected ID %v, got %v", id, userDTO.ID)
	}
	if userDTO.Email.GetEmail() != email.GetEmail() {
		t.Errorf("expected Email %v, got %v", email.GetEmail(), userDTO.Email.GetEmail())
	}
	if len(userDTO.Password.GetPassword()) == 0 {
		t.Error("expected Password to be set")
	}
	if userDTO.IsMentor.GetIsMentor() != isMentor.GetIsMentor() {
		t.Errorf("expected IsMentor %v, got %v", isMentor.GetIsMentor(), userDTO.IsMentor.GetIsMentor())
	}
}

// Для сравнения паролей
func ComparePassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
