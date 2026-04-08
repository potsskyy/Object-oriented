package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/internal/models"
	"todo/internal/repository"
)

func TestRegisterRejectsEmptyCredentials(t *testing.T) {
	store := repository.NewMemoryRepo()
	handler := NewAuthHandler(store)

	payload, err := json.Marshal(models.User{})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("register status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestLoginRejectsWrongPassword(t *testing.T) {
	store := repository.NewMemoryRepo()
	handler := NewAuthHandler(store)

	registerPayload, err := json.Marshal(models.User{
		Username: "admin",
		Password: "1234",
	})
	if err != nil {
		t.Fatalf("marshal register payload: %v", err)
	}

	registerRequest := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(registerPayload))
	registerRequest.Header.Set("Content-Type", "application/json")
	registerResponse := httptest.NewRecorder()
	handler.Register(registerResponse, registerRequest)

	loginPayload, err := json.Marshal(models.User{
		Username: "admin",
		Password: "wrong-password",
	})
	if err != nil {
		t.Fatalf("marshal login payload: %v", err)
	}

	loginRequest := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginPayload))
	loginRequest.Header.Set("Content-Type", "application/json")
	loginResponse := httptest.NewRecorder()
	handler.Login(loginResponse, loginRequest)

	if loginResponse.Code != http.StatusUnauthorized {
		t.Fatalf("login status = %d, want %d", loginResponse.Code, http.StatusUnauthorized)
	}
}
