package pswd

import (
	"golang.org/x/crypto/bcrypt"
	"kubercode-sso/internal/domain/auth/values"
)

func ComparePasswords(userHashedPassword values.Password, incomingPassword string) error {
	result := bcrypt.CompareHashAndPassword(userHashedPassword.GetPassword(), []byte(incomingPassword))
	if result == nil {
		return nil
	}
	return result
}
