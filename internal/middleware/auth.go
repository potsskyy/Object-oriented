package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const UsernameKey contextKey = "username"

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("username")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized: missing username cookie", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UsernameKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUsername(r *http.Request) string {
	username, ok := r.Context().Value(UsernameKey).(string)
	if !ok {
		return ""
	}
	return username
}
