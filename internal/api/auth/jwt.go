package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"banner-service/internal/model"
)

var (
	jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
)

func GenerateToken(tagID int, isAdmin bool) (string, error) {
	claims := &model.Claims{
		TagID:   tagID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, *model.Claims, error) {
	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}
