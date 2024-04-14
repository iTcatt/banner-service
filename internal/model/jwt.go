package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	TagID   int  `json:"tag_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}
