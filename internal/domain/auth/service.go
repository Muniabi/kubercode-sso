package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type Service struct {
	repo        *Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
}

func NewService(repo *Repository, jwtSecret string, tokenExpiry time.Duration) *Service {
	return &Service{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
	}
}

func (s *Service) Register(ctx context.Context, user *User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("[Login] Пользователь не найден: %v", err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("[Login] Сравниваем пароли для пользователя %s", email)
	log.Printf("[Login] Хеш пароля в БД: %s", user.Password)
	log.Printf("[Login] Введенный пароль: %s", password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("[Login] Ошибка сравнения паролей: %v", err)
		return nil, errors.New("invalid credentials")
	}

	// Генерируем access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		log.Printf("[Login] Ошибка генерации access token: %v", err)
		return nil, err
	}

	// Генерируем refresh token
	refreshToken, err := s.generateToken(user)
	if err != nil {
		log.Printf("[Login] Ошибка генерации refresh token: %v", err)
		return nil, err
	}

	// Сохраняем refresh token в базу данных
	token := &Token{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(s.tokenExpiry * 24 * 7), // Refresh token живет 7 дней
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.SaveToken(ctx, token); err != nil {
		log.Printf("[Login] Ошибка сохранения refresh token: %v", err)
		return nil, err
	}

	log.Printf("[Login] Успешный вход пользователя: %s", email)
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			IsMentor: user.IsMentor,
		},
	}, nil
}

func (s *Service) generateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *Service) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
}

func (s *Service) GetUserFromToken(token *jwt.Token) (*User, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		return nil, err
	}

	return s.repo.GetUserByID(context.Background(), userID)
}

// SignUp регистрирует нового пользователя
func (s *Service) SignUp(ctx context.Context, req *SignUpRequest) (*SignUpResponse, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Создаем нового пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:          primitive.NewObjectID(),
		Email:       req.Email,
		Password:    string(hashedPassword),
		IsMentor:    req.IsMentor,
		DeviceToken: req.DeviceToken,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// Генерируем токены
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	// Генерируем refresh token (в реальном приложении здесь может быть другая логика)
	refreshToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &SignUpResponse{
		ID:           user.ID,
		Email:        user.Email,
		IsMentor:     user.IsMentor,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken обновляет токен доступа
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Валидируем refresh token
	token, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Получаем пользователя
	user, err := s.GetUserFromToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Генерируем новый access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

// Logout выполняет выход пользователя
func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	// В реальном приложении здесь можно добавить логику для инвалидации refresh token
	// Например, добавить его в черный список или удалить из базы данных
	return nil
}

// VerifyToken проверяет токен и возвращает информацию о пользователе
func (s *Service) VerifyToken(ctx context.Context, token string) (*UserInfo, error) {
	claims, err := s.ValidateToken(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	mapClaims, ok := claims.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(mapClaims["user_id"].(string))
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &UserInfo{
		ID:       user.ID,
		Email:    user.Email,
		IsMentor: user.IsMentor,
	}, nil
}

// ChangePassword изменяет пароль пользователя
func (s *Service) ChangePassword(ctx context.Context, userID string, oldPassword, newPassword string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidCredentials
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	return s.repo.UpdateUser(ctx, user)
}

// ChangeEmail изменяет email пользователя
func (s *Service) ChangeEmail(ctx context.Context, userID string, newEmail, password string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return ErrInvalidCredentials
	}

	existingUser, err := s.repo.GetUserByEmail(ctx, newEmail)
	if err == nil && existingUser != nil {
		return ErrUserAlreadyExists
	}

	user.Email = newEmail
	return s.repo.UpdateUser(ctx, user)
}

// LogoutFromAllDevices выполняет выход со всех устройств
func (s *Service) LogoutFromAllDevices(ctx context.Context, userID string) error {
	// В реальном приложении здесь можно добавить логику для инвалидации всех refresh токенов
	// Например, удалить все токены из базы данных для данного пользователя
	return nil
} 