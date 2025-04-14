package models

type SignUpRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DeviceToken string `json:"deviceToken" validate:"required"`
	IsMentor    bool   `json:"isMentor"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokensWithUserInfo struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	Account      AccountInfo `json:"account"`
}

type AccountInfo struct {
	Email       string `json:"email"`
	DeviceToken string `json:"deviceToken"`
	IsMentor    bool   `json:"isMentor"`
}

type VerifyTokenResponse struct {
	Status    bool   `json:"status"`
	AccountID string `json:"accountId"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}

type ChangeEmailRequest struct {
	NewEmail string `json:"newEmail" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type OTPRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type OTPResponse struct {
	Status bool `json:"status"`
}

type RestorePasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"newPassword" validate:"required,min=8"`
}