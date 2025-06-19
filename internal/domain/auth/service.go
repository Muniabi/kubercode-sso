package auth

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrTokenRevoked      = errors.New("token has been revoked")
)

type Service struct {
	repo        *Repository
	jwtSecret   []byte
	tokenExpiry time.Duration
	redis       *redis.Client
}

func NewService(repo *Repository, jwtSecret string, tokenExpiry time.Duration, redis *redis.Client) *Service {
	return &Service{
		repo:        repo,
		jwtSecret:   []byte(jwtSecret),
		tokenExpiry: tokenExpiry,
		redis:       redis,
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("[Login] Ошибка сравнения паролей: %v", err)
		return nil, errors.New("invalid credentials")
	}

	// Генерируем access token
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		log.Printf("[Login] Ошибка генерации access token: %v", err)
		return nil, err
	}

	// Генерируем refresh token
	refreshToken, err := s.generateToken(user, true)
	if err != nil {
		log.Printf("[Login] Ошибка генерации refresh token: %v", err)
		return nil, err
	}

	// Сохраняем refresh token в базу данных
	token := &Token{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 30), // Refresh token живет 30 дней
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.SaveToken(ctx, token); err != nil {
		log.Printf("[Login] Ошибка сохранения refresh token: %v", err)
		return nil, err
	}

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

func (s *Service) generateToken(user *User, isRefresh bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour).Unix(), // Access token живет 1 час
	}

	if isRefresh {
		claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix() // Refresh token живет 30 дней
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
func (s *Service) SignUp(ctx context.Context, req *SignUpRequest) (*LoginResponse, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, ErrUserAlreadyExists
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
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(user, true)
	if err != nil {
		return nil, err
	}

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
	accessToken, err := s.generateToken(user, false)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

// Logout выполняет выход пользователя
func (s *Service) Logout(ctx context.Context, userID string) (*gin.H, error) {
	// Получаем токен из контекста
	token := ctx.Value("token").(string)
	if token == "" {
		log.Printf("[Logout] Токен не найден в контексте")
		return nil, errors.New("token not found")
	}

	// Добавляем токен в черный список
	expiration := time.Hour * 24 * 30 // 30 дней
	err := s.redis.Set(ctx, "blacklist:"+token, "revoked", expiration).Err()
	if err != nil {
		log.Printf("[Logout] Ошибка добавления токена в черный список: %v", err)
		return nil, err
	}

	return &gin.H{"message": "Successfully logged out"}, nil
}

// VerifyToken проверяет токен и возвращает информацию о пользователе
func (s *Service) VerifyToken(ctx context.Context, token string) (*UserInfo, error) {
	// Проверяем, не отозван ли токен
	exists, err := s.redis.Exists(ctx, "blacklist:"+token).Result()
	if err != nil {
		log.Printf("[VerifyToken] Ошибка проверки черного списка: %v", err)
		return nil, ErrInvalidToken
	}
	if exists > 0 {
		log.Printf("[VerifyToken] Токен находится в черном списке")
		return nil, ErrTokenRevoked
	}

	// Валидируем токен
	jwtToken, err := s.ValidateToken(token)
	if err != nil {
		log.Printf("[VerifyToken] Ошибка валидации токена: %v", err)
		return nil, ErrInvalidToken
	}

	// Проверяем, что токен не истек
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		log.Printf("[VerifyToken] Неверный формат claims")
		return nil, ErrInvalidToken
	}

	// Проверяем время жизни токена
	exp, ok := claims["exp"].(float64)
	if !ok {
		log.Printf("[VerifyToken] Отсутствует время жизни токена")
		return nil, ErrInvalidToken
	}

	if time.Now().Unix() > int64(exp) {
		log.Printf("[VerifyToken] Токен истек")
		return nil, ErrInvalidToken
	}

	userID, err := primitive.ObjectIDFromHex(claims["user_id"].(string))
	if err != nil {
		log.Printf("[VerifyToken] Неверный формат ID пользователя: %v", err)
		return nil, ErrInvalidToken
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("[VerifyToken] Пользователь не найден: %v", err)
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
func (s *Service) LogoutFromAllDevices(ctx context.Context, userID string) (*gin.H, error) {
	// Получаем текущий токен из контекста
	currentToken := ctx.Value("token").(string)
	if currentToken == "" {
		log.Printf("[LogoutFromAllDevices] Токен не найден в контексте")
		return nil, errors.New("token not found")
	}

	// Добавляем текущий токен в черный список
	expiration := time.Hour * 24 * 30 // 30 дней
	err := s.redis.Set(ctx, "blacklist:"+currentToken, "revoked", expiration).Err()
	if err != nil {
		log.Printf("[LogoutFromAllDevices] Ошибка добавления текущего токена в черный список: %v", err)
		return nil, err
	}

	// Получаем все refresh токены пользователя
	tokens, err := s.repo.GetUserTokens(ctx, userID)
	if err != nil {
		log.Printf("[LogoutFromAllDevices] Ошибка получения токенов пользователя: %v", err)
		return nil, err
	}

	// Добавляем все токены в черный список
	for _, token := range tokens {
		err := s.redis.Set(ctx, "blacklist:"+token.RefreshToken, "revoked", expiration).Err()
		if err != nil {
			log.Printf("[LogoutFromAllDevices] Ошибка добавления токена в черный список: %v", err)
			continue
		}
	}

	// Удаляем все токены из базы данных
	if err := s.repo.DeleteUserTokens(ctx, userID); err != nil {
		log.Printf("[LogoutFromAllDevices] Ошибка удаления токенов из базы данных: %v", err)
		return nil, err
	}

	return &gin.H{"message": "Successfully logged out from all devices"}, nil
} 