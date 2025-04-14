package handlers

import (
	"log"
	"net/http"

	"kubercode/internal/domain/auth"
	"kubercode/internal/domain/models"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

// @Summary     Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.SignUpRequest true "Данные для регистрации"
// @Success     200 {object} models.TokenResponse
// @Failure     400 {object} ErrorResponse
// @Failure     409 {object} ErrorResponse
// @Router      /auth/signup [post]
// @Example     request - {"email": "test@example.com", "password": "password123", "deviceToken": "device123", "isMentor": false}
func (h *AuthHandler) SignUp(c *gin.Context) {
	log.Printf("[SignUp] Получен запрос на регистрацию от %s", c.ClientIP())
	
	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SignUp] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Конвертируем models.SignUpRequest в auth.SignUpRequest
	authReq := &auth.SignUpRequest{
		Email:       req.Email,
		Password:    req.Password,
		DeviceToken: req.DeviceToken,
		IsMentor:    req.IsMentor,
	}

	resp, err := h.service.SignUp(c.Request.Context(), authReq)
	if err != nil {
		switch err {
		case auth.ErrUserAlreadyExists:
			log.Printf("[SignUp] Пользователь уже существует: %s", req.Email)
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		default:
			log.Printf("[SignUp] Внутренняя ошибка сервера: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	log.Printf("[SignUp] Успешная регистрация пользователя: %s", req.Email)
	c.JSON(http.StatusOK, resp)
}

// @Summary     Вход в систему
// @Description Аутентифицирует пользователя и возвращает токены
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.LoginRequest true "Данные для входа"
// @Success     200 {object} models.TokenResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/login [post]
// @Example     request - {"email": "test@example.com", "password": "password123", "deviceToken": "device123"}
func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("[Login] Получен запрос на вход от %s", c.ClientIP())
	
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Login] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		log.Printf("[Login] Ошибка входа: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("[Login] Успешный вход пользователя: %s", req.Email)
	c.JSON(http.StatusOK, resp)
}

// @Summary     Выход из системы
// @Description Выходит пользователя из системы и инвалидирует токены
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} SuccessResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	log.Printf("[Logout] Получен запрос на выход от %s", c.ClientIP())
	
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Printf("[Logout] Токен отсутствует в заголовке")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	if err := h.service.Logout(c.Request.Context(), token); err != nil {
		log.Printf("[Logout] Ошибка при выходе: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	log.Printf("[Logout] Успешный выход пользователя")
	c.Status(http.StatusOK)
}

// @Summary     Обновление токена
// @Description Обновляет access token используя refresh token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} models.TokenResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	log.Printf("[RefreshToken] Получен запрос на обновление токена от %s", c.ClientIP())
	
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Printf("[RefreshToken] Токен отсутствует в заголовке")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), token)
	if err != nil {
		log.Printf("[RefreshToken] Ошибка обновления токена: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to refresh token"})
		return
	}

	log.Printf("[RefreshToken] Токен успешно обновлен")
	c.JSON(http.StatusOK, resp)
}

// @Summary     Изменение пароля
// @Description Изменяет пароль пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body models.ChangePasswordRequest true "Данные для смены пароля"
// @Success     200 {object} SuccessResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/change-password [post]
// @Example     request - {"oldPassword": "oldpass123", "newPassword": "newpass123"}
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	log.Printf("[ChangePassword] Получен запрос на смену пароля от %s", c.ClientIP())
	
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangePassword] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if err := h.service.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		log.Printf("[ChangePassword] Ошибка смены пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ChangePassword] Пароль успешно изменен")
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

// @Summary     Изменение email
// @Description Изменяет email пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request body models.ChangeEmailRequest true "Данные для смены email"
// @Success     200 {object} SuccessResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/change-email [post]
// @Example     request - {"newEmail": "new@example.com", "password": "password123"}
func (h *AuthHandler) ChangeEmail(c *gin.Context) {
	log.Printf("[ChangeEmail] Получен запрос на смену email от %s", c.ClientIP())
	
	var req models.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangeEmail] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if err := h.service.ChangeEmail(c.Request.Context(), userID, req.NewEmail, req.Password); err != nil {
		log.Printf("[ChangeEmail] Ошибка смены email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[ChangeEmail] Email успешно изменен")
	c.JSON(http.StatusOK, gin.H{"message": "Email changed successfully"})
}

// @Summary     Отправка OTP
// @Description Отправляет OTP код на email пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.OTPRequest true "Email для отправки OTP"
// @Success     200 {object} models.OTPResponse
// @Failure     400 {object} ErrorResponse
// @Router      /auth/otp/send [post]
// @Example     request - {"email": "test@example.com"}
func (h *AuthHandler) SendOTP(c *gin.Context) {
	log.Printf("[SendOTP] Получен запрос на отправку OTP от %s", c.ClientIP())
	
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SendOTP] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Тестовая реализация
	response := models.OTPResponse{
		Status: true,
	}

	log.Printf("[SendOTP] OTP успешно отправлен на %s", req.Email)
	c.JSON(http.StatusOK, response)
}

// @Summary     Проверка OTP
// @Description Проверяет OTP код
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.OTPRequest true "OTP код"
// @Success     200 {object} models.OTPResponse
// @Failure     400 {object} ErrorResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/otp/verify [post]
// @Example     request - {"email": "test@example.com", "code": "123456"}
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	log.Printf("[VerifyOTP] Получен запрос на проверку OTP от %s", c.ClientIP())
	
	var req models.OTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[VerifyOTP] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Тестовая реализация
	response := models.OTPResponse{
		Status: true,
	}

	log.Printf("[VerifyOTP] OTP успешно проверен")
	c.JSON(http.StatusOK, response)
}

// @Summary     Восстановление пароля
// @Description Восстанавливает пароль пользователя
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body models.RestorePasswordRequest true "Данные для восстановления пароля"
// @Success     200 {object} SuccessResponse
// @Failure     400 {object} ErrorResponse
// @Router      /auth/restore-password [post]
// @Example     request - {"email": "test@example.com", "newPassword": "newpass123"}
func (h *AuthHandler) RestorePassword(c *gin.Context) {
	log.Printf("[RestorePassword] Получен запрос на восстановление пароля от %s", c.ClientIP())
	
	var req models.RestorePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[RestorePassword] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Тестовая реализация
	response := struct {
		Status bool `json:"status"`
	}{
		Status: true,
	}

	log.Printf("[RestorePassword] Пароль успешно восстановлен")
	c.JSON(http.StatusOK, response)
}

// @Summary     Выход со всех устройств
// @Description Выходит пользователя со всех устройств
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} SuccessResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/logout-all [post]
func (h *AuthHandler) LogoutFromAllDevices(c *gin.Context) {
	log.Printf("[LogoutFromAllDevices] Получен запрос на выход со всех устройств от %s", c.ClientIP())
	
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Printf("[LogoutFromAllDevices] Токен отсутствует в заголовке")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	if err := h.service.LogoutFromAllDevices(c.Request.Context(), token); err != nil {
		log.Printf("[LogoutFromAllDevices] Ошибка выхода со всех устройств: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout from all devices"})
		return
	}

	log.Printf("[LogoutFromAllDevices] Успешный выход со всех устройств")
	c.Status(http.StatusOK)
}

// @Summary     Проверка токена
// @Description Проверяет валидность JWT токена
// @Tags        auth
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} models.VerifyTokenResponse
// @Failure     401 {object} ErrorResponse
// @Router      /auth/verify [post]
func (h *AuthHandler) VerifyToken(c *gin.Context) {
	log.Printf("[VerifyToken] Получен запрос на проверку токена от %s", c.ClientIP())
	
	token := c.GetHeader("Authorization")
	if token == "" {
		log.Printf("[VerifyToken] Токен отсутствует в заголовке")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	resp, err := h.service.VerifyToken(c.Request.Context(), token)
	if err != nil {
		log.Printf("[VerifyToken] Ошибка проверки токена: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	log.Printf("[VerifyToken] Токен успешно проверен")
	c.JSON(http.StatusOK, resp)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
} 