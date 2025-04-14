package values

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestNewPassword_ValidPassword(t *testing.T) {
	passwordStr := "Valid1@password"
	password, err := NewPassword(passwordStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	// Проверка того, что пароль хэширован
	if err := bcrypt.CompareHashAndPassword(password.Password, []byte(passwordStr)); err != nil {
		t.Errorf("expected hashed password to match, got %v", err)
	}
}

func TestNewPassword_InvalidPassword(t *testing.T) {
	invalidPassword := "short"
	_, err := NewPassword(invalidPassword)
	if err == nil {
		t.Errorf("expected error for invalid password, got nil")
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		valid    bool
	}{
		{"Valid1@", true},        // валидный пароль
		{"short", false},         // слишком короткий
		{"NoSpecial1", false},    // нет спецсимвола
		{"nouppercase1@", false}, // нет заглавных букв
		{"NOLOWERCASE1@", false}, // нет строчных букв
		{"NoNumber@", false},     // нет цифры
	}

	for _, tt := range tests {
		if isValidPassword(tt.password) != tt.valid {
			t.Errorf("expected validity for %v to be %v", tt.password, tt.valid)
		}
	}
}

func TestUnmarshalJSON_ValidPassword(t *testing.T) {
	data := []byte(`{"password": "Valid1@password"}`)
	var password Password
	err := json.Unmarshal(data, &password)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	// Проверяем что пароль установлен
	if len(password.Password) == 0 {
		t.Errorf("expected non-empty password")
	}
}

func TestUnmarshalJSON_InvalidPassword(t *testing.T) {
	data := []byte(`{"password": ""}`)
	var password Password
	err := json.Unmarshal(data, &password)
	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

func TestPassword_ToString(t *testing.T) {
	passwordStr := "Valid1@password"
	password, err := NewPassword(passwordStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if password.ToString() == passwordStr {
		t.Errorf("expected hashed password, got plain text")
	}
}

func TestPassword_GetPassword(t *testing.T) {
	passwordStr := "Valid1@password"
	password, err := NewPassword(passwordStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := bcrypt.CompareHashAndPassword(password.GetPassword(), []byte(passwordStr)); err != nil {
		t.Errorf("expected hashed password to match, got %v", err)
	}
}
