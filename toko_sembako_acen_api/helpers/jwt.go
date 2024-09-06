package helpers

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var secretKey = []byte(viper.GetString("JWT_SECRET"))

func CreateToken(userId uuid.UUID, username string, email string) (string, error) {

	// membuat sebuah token yang dimana token tersebut juga berisi sebuah
	// data userid , username , email
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   userId,
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(time.Hour * 4).Unix(),
		},
	)

	// di enkripsi dan juga diconvert menjadi string
	TokenString, err := token.SignedString(secretKey)

	log.Println(token.Claims.GetExpirationTime())

	if err != nil {
		return "", err
	}

	if err != nil {
		log.Println("Error Token Claims Expired Time : " + err.Error())
		return "", err
	}

	return TokenString, nil
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
