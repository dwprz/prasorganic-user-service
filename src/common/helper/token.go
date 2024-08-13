package helper

import (
	"time"

	"github.com/dwprz/prasorganic-user-service/src/infrastructure/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userId string, email string, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":     "prasorganic-auth-service",
		"user_id": userId,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString(config.Conf.Jwt.PrivateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}