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
	repo repository.TaskRepo
}

func NewTaskHandler(repo repository.TaskRepo) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	id, err := h.repo.Add(username, &task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid task", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(username, &task); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = h.repo.Resolve(username, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err = h.repo.Delete(username, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	tasks, err := h.repo.Get(username)
	if err != nil {
		http.Error(w, "error getting tasks", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetArchive(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	tasks, err := h.repo.GetArchive(username)
	if err != nil {
		http.Error(w, "error getting archive", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	username := middleware.GetUsername(r)

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.repo.GetByID(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)
}
