package auth

import (
	"net/http"
	"payment/pkg/errors"
)

type Config struct {
	Token string
}

// AuthMiddleware returns a middleware function that wraps an http.HandlerFunc.
func AuthMiddleware(config *Config) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if len(tokenString) == 0 {
				errors.Error(w, http.StatusUnauthorized, "token not found in header")
				return
			}
			if tokenString != config.Token {
				errors.Error(w, http.StatusUnauthorized, "invalid token")
				return
			}
			next(w, r)
		}
	}
}
