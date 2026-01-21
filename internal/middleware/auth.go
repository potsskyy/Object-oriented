package middleware

import (
	"context"
	"net/http"
)

// contextKeyType используется для ключей в контексте
type contextKeyType string

// UserKey ключ для хранения имени пользователя в контексте
const UserKey contextKeyType = "user_session"

// AuthRequired проверяет наличие cookie и добавляет имя пользователя в контекст
func AuthRequired(next http.Handler) http.Handler {
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

// GetUserFromContext возвращает имя пользователя из контекста запроса
func GetUserFromContext(r *http.Request) string {
	user, ok := r.Context().Value(UserKey).(string)
	if !ok {
		return ""
	}
	return user
}
