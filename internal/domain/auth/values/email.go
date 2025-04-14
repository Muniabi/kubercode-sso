package values

import (
	"encoding/json"
	"fmt"
	"net/mail"
)

type Email struct {
	Email string `json:"email"`
}

func NewEmail(email string) (*Email, error) {
	if !isValidEmail(email) {
		return nil, fmt.Errorf("invalid email: %s", email)
	}

	return &Email{
		Email: email,
	}, nil
}

func (email *Email) UnmarshalJSON(data []byte) error {
	var temp struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if !isValidEmail(temp.Email) {
		return fmt.Errorf("invalid email format: %s", temp.Email)
	}

	email.Email = temp.Email
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (email *Email) ToString() string {
	return fmt.Sprintf("%s", email.Email)
}

func (email *Email) GetEmail() string {
	return email.Email
}
