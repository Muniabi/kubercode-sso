package auth

import (
	"context"
	"errors"
	"time"

	"kubercode/internal/domain/models"
	"kubercode/internal/infrastructure/repository/mongodb"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidToken      = errors.New("invalid token")
)

type Service struct {
	repo            *mongodb.AuthRepository
	accessSecret    string
	refreshSecret   string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewService(
	repo *mongodb.AuthRepository,
	accessSecret string,
	refreshSecret string,
	accessDuration time.Duration,
	refreshDuration time.Duration,
) *Service {
	return &Service{
		repo:            repo,
		accessSecret:    accessSecret,
		refreshSecret:   refreshSecret,
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (s *Service) SignUp(ctx context.Context, req *models.SignUpRequest) (*models.TokensWithUserInfo, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &mongodb.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		DeviceToken: req.DeviceToken,
		IsMentor:    req.IsMentor,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Генерируем токены
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh token
	if err := s.repo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshDuration)); err != nil {
		return nil, err
	}

	return &models.TokensWithUserInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account: models.AccountInfo{
			Email:       user.Email,
			DeviceToken: user.DeviceToken,
			IsMentor:    user.IsMentor,
		},
	}, nil
}

func (s *Service) Login(ctx context.Context, req *models.LoginRequest) (*models.TokensWithUserInfo, error) {
	// Находим пользователя
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidPassword
	}

	// Обновляем device token
	user.DeviceToken = req.DeviceToken
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	// Генерируем токены
	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, err
	}

	// Сохраняем refresh token
	if err := s.repo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(s.refreshDuration)); err != nil {
		return nil, err
	}

	return &models.TokensWithUserInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Account: models.AccountInfo{
			Email:       user.Email,
			DeviceToken: user.DeviceToken,
			IsMentor:    user.IsMentor,
		},
	}, nil
}

func (s *Service) VerifyToken(ctx context.Context, token string) (*models.VerifyTokenResponse, error) {
	claims, err := s.parseToken(token, s.accessSecret)
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex((*claims)["sub"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &models.VerifyTokenResponse{
		Status:    true,
		AccountID: user.ID.Hex(),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	userID, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if userID == nil {
		return ErrInvalidToken
	}

	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	userID, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	if userID == nil {
		return nil, ErrInvalidToken
	}

	// Генерируем новые токены
	accessToken, newRefreshToken, err := s.generateTokens(*userID)
	if err != nil {
		return nil, err
	}

	// Удаляем старый refresh token
	if err := s.repo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Сохраняем новый refresh token
	if err := s.repo.StoreRefreshToken(ctx, *userID, newRefreshToken, time.Now().Add(s.refreshDuration)); err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	// Находим пользователя
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ErrInvalidToken
	}

	user, err := s.repo.FindUserByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Проверяем старый пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Хешируем новый пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Обновляем пароль
	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) ChangeEmail(ctx context.Context, userID string, req *models.ChangeEmailRequest) error {
	// Проверяем, не занят ли email
	existingUser, err := s.repo.FindUserByEmail(ctx, req.NewEmail)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return ErrUserAlreadyExists
	}

	// Находим пользователя
	user, err := s.repo.FindUserByEmail(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Обновляем email
	user.Email = req.NewEmail
	return s.repo.UpdateUser(ctx, user)
}

func (s *Service) LogoutFromAllDevices(ctx context.Context, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteAllUserRefreshTokens(ctx, id)
}

func (s *Service) generateTokens(userID primitive.ObjectID) (string, string, error) {
	// Генерируем access token
	accessClaims := jwt.MapClaims{
		"sub": userID.Hex(),
		"exp": time.Now().Add(s.accessDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", "", err
	}

	// Генерируем refresh token
	refreshClaims := jwt.MapClaims{
		"sub": userID.Hex(),
		"exp": time.Now().Add(s.refreshDuration).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *Service) parseToken(tokenString, secret string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, ErrInvalidToken
} 