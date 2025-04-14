package values

import (
	"encoding/json"
	"testing"
)

func TestNewEmail_ValidEmail(t *testing.T) {
	emailStr := "example@example.com"
	email, err := NewEmail(emailStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if email.Email != emailStr {
		t.Errorf("expected %v, got %v", emailStr, email.Email)
	}
}

func TestNewEmail_InvalidEmail(t *testing.T) {
	invalidEmail := "invalid-email"
	_, err := NewEmail(invalidEmail)
	if err == nil {
		t.Errorf("expected error for invalid email, got nil")
	}
}

func TestUnmarshalJSON_ValidEmail(t *testing.T) {
	data := []byte(`{"email": "test@example.com"}`)
	var email Email
	err := json.Unmarshal(data, &email)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if email.Email != "test@example.com" {
		t.Errorf("expected test@example.com, got %v", email.Email)
	}
}

func TestUnmarshalJSON_InvalidEmail(t *testing.T) {
	data := []byte(`{"email": "invalid-email"}`)
	var email Email
	err := json.Unmarshal(data, &email)
	if err == nil {
		t.Errorf("expected error for invalid email, got nil")
	}
}

func TestEmail_ToString(t *testing.T) {
	emailStr := "example@example.com"
	email, err := NewEmail(emailStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if email.ToString() != emailStr {
		t.Errorf("expected %v, got %v", emailStr, email.ToString())
	}
}

func TestEmail_GetEmail(t *testing.T) {
	emailStr := "example@example.com"
	email, err := NewEmail(emailStr)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if email.GetEmail() != emailStr {
		t.Errorf("expected %v, got %v", emailStr, email.GetEmail())
	}
}
