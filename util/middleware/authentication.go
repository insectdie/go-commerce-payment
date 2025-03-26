package middleware

import (
	"codebase-service/helper"
	"context"
	"net/http"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(userIDKey).(string)
}

func SetRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

func GetRole(ctx context.Context) string {
	return ctx.Value(roleKey).(string)
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		if token == "" {
			helper.HandleResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		token = token[len("Bearer "):]
		payload, err := VerifyToken(token)
		if err != nil {
			helper.HandleResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		ctx = SetUserID(ctx, payload.UserID)
		ctx = SetRole(ctx, payload.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userId := r.Header.Get("X-USER-ID")

		if userId == "" {
			helper.HandleResponse(w, http.StatusUnauthorized, "Unauthorized", nil)
			return
		}

		ctx = SetUserID(ctx, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
