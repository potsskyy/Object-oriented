package middleware

import (
	"context"
	"net/http"

	"todo/internal/repository"
	"todo/internal/session"
)

type contextKeyType string

const UserKey contextKeyType = "authenticated_user"

func WithAuth(store repository.SessionRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(session.CookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, "unauthorized: missing session cookie", http.StatusUnauthorized)
				return
			}

			username, err := store.GetUsernameBySession(cookie.Value)
			if err != nil || username == "" {
				http.Error(w, "unauthorized: invalid session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUsername(r *http.Request) string {
	user, ok := r.Context().Value(UserKey).(string)
	if !ok {
		return ""
	}

	return user
}
