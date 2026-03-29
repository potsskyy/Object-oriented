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

type TaskHandler struct {
	store repository.TaskRepo
}

func NewTaskHandler(store repository.TaskRepo) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	id, err := h.store.Add(user, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	if err := h.store.Update(user, &task); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.store.Resolve(user, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.store.Delete(user, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.store.Get(user)
	if err != nil {
		http.Error(w, "error getting tasks", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.store.GetByID(user, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(task)
}

func (h *TaskHandler) GetArchive(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUsername(r)
	if user == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	tasks, err := h.store.GetArchive(user)
	if err != nil {
		http.Error(w, "error getting archive", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
}
