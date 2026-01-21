package models

import "time"

// Todo представляет задачу пользователя
type Todo struct {
	ID         int64      `json:"id"`
	User       string     `json:"user"`
	Headline   string     `json:"headline"`
	Details    string     `json:"details"`
	Completed  bool       `json:"completed"`
	Archived   bool       `json:"archived"`
	CreatedOn  time.Time  `json:"created_on"`
	UpdatedOn  time.Time  `json:"updated_on"`
	FinishedOn *time.Time `json:"finished_on,omitempty"`
}
