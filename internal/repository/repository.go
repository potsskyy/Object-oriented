package repository

import "todo/internal/models"

type UserRepo interface {
	Register(user *models.User) error
	GetUser(username string) (*models.User, error)
}

type TaskRepo interface {
	Add(username string, task *models.Task) (int64, error)
	Get(username string) ([]*models.Task, error)
	GetArchive(username string) ([]*models.Task, error)
	GetByID(username string, id int64) (*models.Task, error)
	Update(username string, updated *models.Task) error
	Resolve(username string, id int64) error
	Delete(username string, id int64) error
}
