package router

import (
	"net/http"

	"kubercode-sso/internal/infrastructure/http/handlers"
	"kubercode-sso/internal/infrastructure/http/middleware"
)

func NewRouter(authHandler *handlers.AuthHandler) http.Handler {
	mux := http.NewServeMux()

	// Публичные эндпоинты
	mux.HandleFunc("/api/v1/auth/signup", authHandler.SignUp)
	mux.HandleFunc("/api/v1/auth/login", authHandler.Login)
	mux.HandleFunc("/api/v1/auth/verify-token", authHandler.VerifyToken)
	mux.HandleFunc("/api/v1/auth/otp/send", authHandler.SendOTP)
	mux.HandleFunc("/api/v1/auth/otp/verify", authHandler.VerifyOTP)
	mux.HandleFunc("/api/v1/auth/password/restore", authHandler.RestorePassword)

	// Защищенные эндпоинты
	protected := http.NewServeMux()
	protected.HandleFunc("/api/v1/auth/logout", authHandler.Logout)
	protected.HandleFunc("/api/v1/auth/refresh-token", authHandler.RefreshToken)
	protected.HandleFunc("/api/v1/auth/password", authHandler.ChangePassword)
	protected.HandleFunc("/api/v1/auth/email", authHandler.ChangeEmail)
	protected.HandleFunc("/api/v1/auth/logout-all-devices", authHandler.LogoutFromAllDevices)

	// Применяем middleware к защищенным эндпоинтам
	protectedWithAuth := middleware.AuthMiddleware(protected)
	mux.Handle("/api/v1/auth/", protectedWithAuth)

	// Применяем CORS middleware ко всем эндпоинтам
	handler := middleware.CORSMiddleware(mux)

	return handler
} 