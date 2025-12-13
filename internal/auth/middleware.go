package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/api"
)

const tokenContextKey contextKey = "token"

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.Unauthorized(w, ErrInvalidToken.Error())
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			api.Unauthorized(w, ErrInvalidToken.Error())
			return
		}

		tokenString := parts[1]

		data, err := h.service.ValidateToken(r.Context(), ValidateTokenParams{Token: tokenString})
		if err != nil {
			api.Unauthorized(w, ErrInvalidToken.Error())
			return
		}

		userID, ok := data["id"].(string)
		if !ok {
			api.Unauthorized(w, ErrInvalidToken.Error())
			return
		}

		role, ok := data["role"].(string)
		if !ok {
			api.Unauthorized(w, ErrInvalidToken.Error())
			return
		}

		// Store token and user info in context
		ctx := context.WithValue(r.Context(), tokenContextKey, tokenString)
		ctx = SetUserContext(ctx, userID, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := GetRole(r.Context())
		if !ok || role != "admin" {
			api.Forbidden(w, "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
