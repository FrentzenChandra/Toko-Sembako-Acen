package helpers

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

var secretKey = []byte(viper.GetString("JWT_SECRET"))

func CreateToken(userId uuid.UUID, username string, email string) (refreshTokenBack string, accessTokenBack string, err error) {

	// membuat sebuah token yang dimana token tersebut juga berisi sebuah
	// data userid , username , email
	Accesstoken := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId":   userId,
			"username": username,
			"email":    email,
			"exp":      time.Now().Add(time.Minute * 15).Unix(),
		},
	)

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	refreshTokenString, err := refreshToken.SignedString(secretKey)

	if err != nil {
		return "", "", err
	}

	// di enkripsi dan juga diconvert menjadi string
	AccessTokenString, err := Accesstoken.SignedString(secretKey)

	if err != nil {
		return "", "", err
	}

	return refreshTokenString, AccessTokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := GetTokenFromHeader(tokenString)

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid Token")
	}

	return nil
}

func GetTokenFromHeader(tokenString string) (*jwt.Token, error) {

	if !strings.Contains(tokenString, "Bearer") {
		return nil, errors.New("Token String Invalid")
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, err
}

func GetTokenFromString(tokenString string) (*jwt.Token, error) {

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	return token, err
}

func ExtractTokenData(token *jwt.Token) (map[string]interface{}, error) {

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func CreateAnotherAccessToken(RefreshTokenString string, accessTokenString string) (refreshTokenBack string, accessTokenBack string, err error) {

	Refreshtoken, err := jwt.Parse(RefreshTokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})

	if claims, ok := Refreshtoken.Claims.(jwt.MapClaims); ok {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			accesstoken, err := GetTokenFromString(accessTokenString)

			if err != nil {
				log.Println("Error When getting token from access Token string : " + err.Error())
				return "", "", err
			}

			data, err := ExtractTokenData(accesstoken)

			userUUID := uuid.MustParse(data["userId"].(string))
			userName := data["username"].(string)
			email := data["email"].(string)

			if userUUID == uuid.Nil || userName == "" || email == "" {
				return "", "", errors.New("Invalid Access Token")
			}

			if err != nil {
				return "", "", err
			}

			refreshToken, accessToken, err := CreateToken(userUUID, userName, email)

			if err != nil {
				return "", "", err
			}

			return refreshToken, accessToken, nil
		}

		return "", "", errors.New("Invalid Access Token")
	}

	return "", "", errors.New("Invalid Refresh Token")
}
