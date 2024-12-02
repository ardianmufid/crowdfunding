package utils

import (
	"crowdfunding/config"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId int) (string, error) {
	claims := jwt.MapClaims{
		"id": userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.Cfg.App.Encryption.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (*jwt.Token, error) {
	tokens, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return tokens, err
	}

	// claims, ok := tokens.Claims.(jwt.MapClaims)
	// if ok && tokens.Valid {
	// 	id := fmt.Sprintf("%v", claims["id"])
	// 	return int(id), nil
	// }

	return tokens, nil
}
