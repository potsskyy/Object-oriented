package repository

import (
	"errors"
	"sync"
	"time"

	"todo/internal/models"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrInvalidLogin = errors.New("invalid username or password")
	ErrTaskNotFound = errors.New("task not found")
)

type MemoryRepo struct {
	mu         sync.RWMutex
	users      map[string]*models.User   // username -> user
	tasks      map[string][]*models.Task // username -> tasks
	nextTaskID int64
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users:      make(map[string]*models.User),
		tasks:      make(map[string][]*models.Task),
		nextTaskID: 1,
	}
}

func (r *MemoryRepo) Register(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Username]; exists {
		return ErrUserExists
	}
	r.users[user.Username] = &models.User{
		Username: user.Username,
		Password: user.Password,
	}
	return nil
}

func (r *MemoryRepo) GetUser(username string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[username]
	if !ok {
		return nil, ErrInvalidLogin
	}
	return user, nil
}

func (r *MemoryRepo) Add(username string, task *models.Task) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[username]; !ok {
		return 0, ErrInvalidLogin
	}

	now := time.Now()
	task.ID = r.nextTaskID
	task.Owner = username // В in-memory избыточно
	task.CreatedAt = now
	task.UpdatedAt = now
	r.nextTaskID++

	r.tasks[username] = append(r.tasks[username], task)
	return task.ID, nil
}

func (r *MemoryRepo) Get(username string) ([]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var active []*models.Task
	for _, task := range r.tasks[username] {
		if task.IsResolved || task.IsArchived {
			continue
		}
		active = append(active, task)
	}
	return active, nil
}

func (r *MemoryRepo) GetArchive(username string) ([]*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var archived []*models.Task
	for _, task := range r.tasks[username] {
		if !task.IsResolved && !task.IsArchived {
			continue
		}
		archived = append(archived, task)
	}
	return archived, nil
}

func (r *MemoryRepo) GetByID(username string, id int64) (*models.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks[username] {
		if task.ID != id {
			continue
		}
		return task, nil
	}
	return nil, ErrTaskNotFound
}

func (r *MemoryRepo) Update(username string, updated *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, task := range r.tasks[username] {
		if task.ID != updated.ID {
			continue
		}
		task.Title = updated.Title
		task.Description = updated.Description
		task.UpdatedAt = time.Now()
		return nil
	}
	return ErrTaskNotFound
}

func (r *MemoryRepo) Resolve(username string, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, task := range r.tasks[username] {
		if task.ID != id {
			continue
		}
		now := time.Now()
		task.IsResolved = true
		task.IsArchived = true
		task.ResolvedAt = &now
		task.UpdatedAt = now
		return nil
	}
	return ErrTaskNotFound
}

func (r *MemoryRepo) Delete(username string, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, task := range r.tasks[username] {
		if task.ID != id {
			continue
		}
		task.IsArchived = true
		task.UpdatedAt = time.Now()
		return nil
	}
	return ErrTaskNotFound
}
