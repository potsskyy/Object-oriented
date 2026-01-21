package repository

import (
	"errors"
	"sync"
	"time"

	"todo/internal/models"
)

var (
	ErrAccountExists   = errors.New("account already exists")
	ErrInvalidAccount  = errors.New("invalid user id or secret")
	ErrTodoNotFound    = errors.New("todo not found")
)

// MemoryStore хранит пользователей и задачи в памяти
type MemoryStore struct {
	mu        sync.RWMutex
	accounts  map[string]*models.Account   // user_id -> account
	todos     map[string][]*models.Todo    // user_id -> todos
	nextTodoID int64
}

// NewMemoryStore создаёт новый in-memory репозиторий
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		accounts:  make(map[string]*models.Account),
		todos:     make(map[string][]*models.Todo),
		nextTodoID: 1,
	}
}

// Register добавляет нового пользователя
func (m *MemoryStore) Register(account *models.Account) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.accounts[account.UserID]; exists {
		return ErrAccountExists
	}
	m.accounts[account.UserID] = &models.Account{
		UserID: account.UserID,
		Secret: account.Secret,
	}
	return nil
}

// GetAccount возвращает пользователя по user_id
func (m *MemoryStore) GetAccount(userID string) (*models.Account, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	acc, ok := m.accounts[userID]
	if !ok {
		return nil, ErrInvalidAccount
	}
	return acc, nil
}

// AddTodo добавляет новую задачу для пользователя
func (m *MemoryStore) AddTodo(userID string, todo *models.Todo) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.accounts[userID]; !ok {
		return 0, ErrInvalidAccount
	}

	now := time.Now()
	todo.ID = m.nextTodoID
	todo.User = userID
	todo.CreatedOn = now
	todo.UpdatedOn = now
	m.nextTodoID++

	m.todos[userID] = append(m.todos[userID], todo)
	return todo.ID, nil
}

// GetTodos возвращает все активные задачи пользователя
func (m *MemoryStore) GetTodos(userID string) ([]*models.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active []*models.Todo
	for _, t := range m.todos[userID] {
		if t.Completed || t.Archived {
			continue
		}
		active = append(active, t)
	}
	return active, nil
}

// GetArchivedTodos возвращает архивные задачи пользователя
func (m *MemoryStore) GetArchivedTodos(userID string) ([]*models.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var archived []*models.Todo
	for _, t := range m.todos[userID] {
		if !t.Completed && !t.Archived {
			continue
		}
		archived = append(archived, t)
	}
	return archived, nil
}

// GetTodoByID возвращает задачу по ID
func (m *MemoryStore) GetTodoByID(userID string, id int64) (*models.Todo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, t := range m.todos[userID] {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, ErrTodoNotFound
}

// UpdateTodo обновляет задачу пользователя
func (m *MemoryStore) UpdateTodo(userID string, updated *models.Todo) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.todos[userID] {
		if t.ID == updated.ID {
			t.Headline = updated.Headline
			t.Details = updated.Details
			t.UpdatedOn = time.Now()
			return nil
		}
	}
	return ErrTodoNotFound
}

// CompleteTodo помечает задачу как выполненную
func (m *MemoryStore) CompleteTodo(userID string, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.todos[userID] {
		if t.ID == id {
			now := time.Now()
			t.Completed = true
			t.Archived = true
			t.FinishedOn = &now
			t.UpdatedOn = now
			return nil
		}
	}
	return ErrTodoNotFound
}

// RemoveTodo архивирует задачу
func (m *MemoryStore) RemoveTodo(userID string, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, t := range m.todos[userID] {
		if t.ID == id {
			t.Archived = true
			t.UpdatedOn = time.Now()
			return nil
		}
	}
	return ErrTodoNotFound
}
