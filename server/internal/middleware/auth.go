package middleware

import (
	"net/http"
)

type authService interface {
	ValidateToken(token string) (string, error)
}

type AuthMiddleware struct {
	authService authService
}

func NewAuthMiddleware(authService authService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (a *AuthMiddleware) Authenticate(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Token is required", http.StatusUnauthorized)
			return
		}

		user, err := a.authService.ValidateToken(token)
		if err != nil || user == "" {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
