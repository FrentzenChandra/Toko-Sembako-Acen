package helpers

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("tokosembakoapijwtsecret12345")

func CreateToken(userId uuid.UUID, username string, email string) (string, error) {

	// membuat sebuah token yang dimana token tersebut juga berisi sebuah
	// data userid , username , email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   userId,
			"username": username,
			"email":    email,
		})

	// di enkripsi dan juga diconvert menjadi string
	signedToken, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString string) error {
	token, err := GetToken(tokenString)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid Token")
	}

	return nil
}

func GetToken(tokenString string) (*jwt.Token, error) {

	if !strings.Contains(tokenString, "Bearer") {
		return nil, errors.New("Token String Invalid")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, err
}

func ExtractTokenData(token *jwt.Token) (map[string]interface{}, error) {

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}
