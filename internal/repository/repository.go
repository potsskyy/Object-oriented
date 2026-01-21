package repository

import "todo/internal/models"

// AccountRepo интерфейс для работы с пользователями
type AccountRepo interface {
	Register(account *models.Account) error
	GetAccount(userID string) (*models.Account, error)
}

// TodoRepo интерфейс для работы с задачами
type TodoRepo interface {
	AddTodo(userID string, todo *models.Todo) (int64, error)
	GetTodos(userID string) ([]*models.Todo, error)
	GetArchivedTodos(userID string) ([]*models.Todo, error)
	GetTodoByID(userID string, id int64) (*models.Todo, error)
	UpdateTodo(userID string, updated *models.Todo) error
	CompleteTodo(userID string, id int64) error
	RemoveTodo(userID string, id int64) error
}
