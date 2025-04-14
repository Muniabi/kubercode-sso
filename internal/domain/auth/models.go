package auth

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User представляет модель пользователя
type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string            `bson:"email" json:"email"`
	Password    string            `bson:"password" json:"-"`
	IsMentor    bool              `bson:"is_mentor" json:"is_mentor"`
	DeviceToken string            `bson:"device_token" json:"device_token"`
	CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// SignUpRequest представляет запрос на регистрацию
type SignUpRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	IsMentor    bool   `json:"is_mentor"`
	DeviceToken string `json:"deviceToken" binding:"required"`
}

// SignUpResponse представляет ответ на регистрацию
type SignUpResponse struct {
	ID           primitive.ObjectID `json:"id"`
	Email        string            `json:"email"`
	IsMentor     bool              `json:"is_mentor"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserInfo представляет информацию о пользователе
type UserInfo struct {
	ID       primitive.ObjectID `json:"id"`
	Email    string            `json:"email"`
	IsMentor bool              `json:"is_mentor"`
}

// LoginResponse представляет ответ на вход
type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	User         UserInfo  `json:"user"`
}

// RefreshTokenRequest представляет запрос на обновление токена
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse представляет ответ на обновление токена
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// LogoutRequest представляет запрос на выход
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ChangePasswordRequest представляет запрос на изменение пароля
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ChangeEmailRequest представляет запрос на изменение email
type ChangeEmailRequest struct {
	NewEmail string `json:"newEmail" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// VerifyOTPRequest представляет запрос на проверку OTP
type VerifyOTPRequest struct {
	OTP string `json:"otp" binding:"required,len=6"`
} 