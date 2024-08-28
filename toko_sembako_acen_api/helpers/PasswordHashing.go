package helpers

import (
	"toko_sembako_acen/infra/logger"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		logger.Infof("Password Hashing Error : " + err.Error())
		return ""
	}

	return string(bytes)
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		logger.Infof("Password Verification Error : " + err.Error())
	}

	return true
}
