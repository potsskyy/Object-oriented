package repository

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"todo/internal/models"
	"todo/internal/session"
)

var (
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidUser    = errors.New("invalid username or password")
	ErrInvalidSession = errors.New("invalid session")
	ErrTaskNotFound   = errors.New("task not found")
)

type userSession struct {
	Username  string
	ExpiresAt time.Time
}

type MemoryRepo struct {
	mu         sync.RWMutex
	users      map[string]*models.User
	sessions   map[string]userSession
	tasks      map[string][]*models.Task
	nextTaskID int64
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		users:      make(map[string]*models.User),
		sessions:   make(map[string]userSession),
		tasks:      make(map[string][]*models.Task),
		nextTaskID: 1,
	}
}

func (m *MemoryRepo) Register(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.Username]; exists {
		return ErrUserExists
	}

	m.users[user.Username] = &models.User{
		Username: user.Username,
		Password: user.Password,
	}

	return nil
}

func (m *MemoryRepo) GetUser(username string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, ok := m.users[username]
	if !ok {
		return nil, ErrInvalidUser
	}

	return user, nil
}

func (m *MemoryRepo) CreateSession(username string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[username]; !ok {
		return "", ErrInvalidUser
	}

	sessionID, err := newSessionID()
	if err != nil {
		return "", err
	}

	m.sessions[sessionID] = userSession{
		Username:  username,
		ExpiresAt: time.Now().Add(session.TTL),
	}

	return sessionID, nil
}

func (m *MemoryRepo) GetUsernameBySession(sessionID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	savedSession, ok := m.sessions[sessionID]
	if !ok {
		return "", ErrInvalidSession
	}

	if time.Now().After(savedSession.ExpiresAt) {
		delete(m.sessions, sessionID)
		return "", ErrInvalidSession
	}

	return savedSession.Username, nil
}

func (m *MemoryRepo) DeleteSession(sessionID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.sessions, sessionID)
}

func (m *MemoryRepo) Add(username string, task *models.Task) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.users[username]; !ok {
		return 0, ErrInvalidUser
	}

	now := time.Now()
	task.ID = m.nextTaskID
	task.User = username
	task.CreatedOn = now
	task.UpdatedOn = now
	m.nextTaskID++

	m.tasks[username] = append(m.tasks[username], task)

	return task.ID, nil
}

func (m *MemoryRepo) Get(username string) ([]*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	active := make([]*models.Task, 0, len(m.tasks[username]))
	for _, task := range m.tasks[username] {
		if task.Completed || task.Archived {
			continue
		}
		active = append(active, task)
	}

	return active, nil
}

func (m *MemoryRepo) GetArchive(username string) ([]*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	archived := make([]*models.Task, 0, len(m.tasks[username]))
	for _, task := range m.tasks[username] {
		if !task.Archived {
			continue
		}
		archived = append(archived, task)
	}

	return archived, nil
}

func (m *MemoryRepo) GetByID(username string, id int64) (*models.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, task := range m.tasks[username] {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, ErrTaskNotFound
}

func (m *MemoryRepo) Update(username string, updated *models.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, task := range m.tasks[username] {
		if task.ID != updated.ID {
			continue
		}

		task.Headline = updated.Headline
		task.Details = updated.Details
		task.UpdatedOn = time.Now()
		return nil
	}

	return ErrTaskNotFound
}

func newSessionID() (string, error) {
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes), nil
}

func (m *MemoryRepo) Resolve(username string, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, task := range m.tasks[username] {
		if task.ID != id {
			continue
		}

		now := time.Now()
		task.Completed = true
		task.Archived = true
		task.FinishedOn = &now
		task.UpdatedOn = now
		return nil
	}

	return ErrTaskNotFound
}

func (m *MemoryRepo) Delete(username string, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, task := range m.tasks[username] {
		if task.ID != id {
			continue
		}

		task.Archived = true
		task.UpdatedOn = time.Now()
		return nil
	}

	return ErrTaskNotFound
}
