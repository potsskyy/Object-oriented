package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"todo/internal/middleware"
	"todo/internal/models"
	"todo/internal/repository"
)

// TodoHandler управляет задачами пользователя
type TodoHandler struct {
	store repository.TaskRepo
}

// NewTodoHandler создаёт новый обработчик задач
func NewTodoHandler(store repository.TaskRepo) *TodoHandler {
	return &TodoHandler{store: store}
}

// Create добавляет новую задачу
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	var todo models.Task
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	id, err := h.store.Add(user, &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

// Modify обновляет существующую задачу
func (h *TodoHandler) Modify(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	var todo models.Task
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	if err := h.store.Update(user, &todo); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Complete помечает задачу как выполненную
func (h *TodoHandler) Complete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = h.store.Resolve(user, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Remove удаляет задачу
func (h *TodoHandler) Remove(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = h.store.Delete(user, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// List возвращает все текущие задачи пользователя
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	todos, err := h.store.Get(user)
	if err != nil {
		http.Error(w, "error getting tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

// Archive возвращает архивные задачи
func (h *TodoHandler) Archive(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	todos, err := h.store.GetArchive(user)
	if err != nil {
		http.Error(w, "error getting archive", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

// GetOne возвращает задачу по ID
func (h *TodoHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	todo, err := h.store.GetByID(user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(todo)
}
