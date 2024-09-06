package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Println("Password Hashing Error : " + err.Error())
		return ""
	}

	return string(bytes)
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(hash, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))


	if err != nil {
		log.Println("Password Verification Error : " + err.Error())
		return false
	}


	return true
}
