package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"

	"banner-service/internal/api/auth"
	"banner-service/internal/model"
	"banner-service/internal/service"
)

var (
	ErrValidationFailed = errors.New("validation failed")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNoPermission     = errors.New("permission denied")
)

func errorsMiddleware(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		switch {
		case err == nil:
			return
		case errors.Is(err, ErrValidationFailed) || errors.Is(err, service.ErrAlreadyExists):
			_ = sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		case errors.Is(err, ErrUnauthorized):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrNoPermission):
			w.WriteHeader(http.StatusForbidden)
		case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
			w.WriteHeader(http.StatusUnauthorized)
		case errors.Is(err, jwt.ErrTokenMalformed):
			w.WriteHeader(http.StatusBadRequest)
		default:
			_ = sendJSONResponse(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}
	}
}

func authMiddleware(_ http.ResponseWriter, r *http.Request) (*model.Claims, error) {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		return nil, ErrUnauthorized
	}
	token, claims, err := auth.VerifyToken(tokenString)
	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

func sendJSONResponse(w http.ResponseWriter, response interface{}, status int) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	log.Printf("sending response")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return err
	}
	return nil
}
