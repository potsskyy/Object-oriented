package middleware

import (
	"context"
	"net/http"
)

type contextKeyType string

const UserKey contextKeyType = "user_session"

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_session")
		if err != nil || cookie.Value == "" {
			http.Error(w, "unauthorized: missing user cookie", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUsername(r *http.Request) string {
	user, ok := r.Context().Value(UserKey).(string)
	if !ok {
		return ""
	}

	return user
}
