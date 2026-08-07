package auth

import (
	"fmt"
	"os"
	"time"
	"url-shortener/internal/helper"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func GenerateToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)

}

func ValidateToken(tokenString string) (int, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, helper.ErrInvalidToken
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, helper.ErrInvalidToken
	}

	return int(userIDFloat), nil

}
