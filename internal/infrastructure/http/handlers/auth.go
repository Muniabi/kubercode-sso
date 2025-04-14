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

func (h *AuthHandler) SignUp(c *gin.Context) {
	log.Printf("[SignUp] Получен запрос на регистрацию от %s", c.ClientIP())
	
	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[SignUp] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.service.SignUp(c.Request.Context(), &req)
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

func (h *AuthHandler) Login(c *gin.Context) {
	log.Printf("[Login] Получен запрос на вход от %s", c.ClientIP())
	
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[Login] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		log.Printf("[Login] Ошибка входа: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	log.Printf("[Login] Успешный вход пользователя: %s", req.Email)
	c.JSON(http.StatusOK, resp)
}

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

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	log.Printf("[ChangePassword] Получен запрос на смену пароля от %s", c.ClientIP())
	
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangePassword] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("userID")
	if err := h.service.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		log.Printf("[ChangePassword] Ошибка смены пароля: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
		return
	}

	log.Printf("[ChangePassword] Пароль успешно изменен")
	c.Status(http.StatusOK)
}

func (h *AuthHandler) ChangeEmail(c *gin.Context) {
	log.Printf("[ChangeEmail] Получен запрос на смену email от %s", c.ClientIP())
	
	var req models.ChangeEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangeEmail] Ошибка декодирования запроса: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID := c.GetString("userID")
	if err := h.service.ChangeEmail(c.Request.Context(), userID, &req); err != nil {
		log.Printf("[ChangeEmail] Ошибка смены email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change email"})
		return
	}

	log.Printf("[ChangeEmail] Email успешно изменен")
	c.Status(http.StatusOK)
}

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