package values

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"

	//"golang.org/x/crypto/bcrypt"
	"unicode"
)

type Password struct {
	Password []byte `json:"password"`
}

func NewPassword(password string) (*Password, error) {
	if !isValidPassword(password) {
		return nil, fmt.Errorf("invalid password: %s", password)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Password{
		Password: hashedPassword,
	}, nil
}

func isValidPassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 7 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func (password *Password) UnmarshalJSON(data []byte) error {
	var temp struct {
		Password []byte `json:"password"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	password.Password = temp.Password
	return nil
}

func (password *Password) ToString() string {
	return fmt.Sprintf("%s", password.Password)
}

func (password *Password) GetPassword() []byte {
	return password.Password
}
