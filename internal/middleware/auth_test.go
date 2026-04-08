package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/models"
	"todo/internal/repository"
	"todo/internal/session"
)

func TestWithAuthRejectsForgedCookie(t *testing.T) {
	store := repository.NewMemoryRepo()

	protected := WithAuth(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  session.CookieName,
		Value: "admin",
	})
	rr := httptest.NewRecorder()

	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("middleware status = %d, want %d", rr.Code, http.StatusUnauthorized)
	}
}

func TestWithAuthPassesAuthenticatedUsernameToContext(t *testing.T) {
	store := repository.NewMemoryRepo()

	if err := store.Register(&models.User{
		Username: "student",
		Password: "hashed-password",
	}); err != nil {
		t.Fatalf("register user: %v", err)
	}

	sessionID, err := store.CreateSession("student")
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	protected := WithAuth(store)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username := GetUsername(r); username != "student" {
			t.Fatalf("username from context = %q, want %q", username, "student")
		}

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/get", nil)
	req.AddCookie(&http.Cookie{
		Name:  session.CookieName,
		Value: sessionID,
	})
	rr := httptest.NewRecorder()

	protected.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("middleware status = %d, want %d", rr.Code, http.StatusOK)
	}
}
