package handler

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"todo/internal/models"
	"todo/internal/repository"
)

// UserHandler отвечает за регистрацию и авторизацию пользователей
type UserHandler struct {
	store repository.UserRepo
}

// NewUserHandler создаёт новый обработчик пользователей
func NewUserHandler(store repository.UserRepo) *UserHandler {
	return &UserHandler{store: store}
}

// SignUp регистрирует нового пользователя
func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)

	if err = h.store.Register(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// SignIn выполняет вход пользователя
func (h *UserHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var loginReq models.User
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	savedUser, err := h.store.GetUser(loginReq.Username)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(savedUser.Password), []byte(loginReq.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_session",
		Value:    loginReq.Username,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}
