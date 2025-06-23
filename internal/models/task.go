package models

import "time"

type Task struct {
	ID          int64      `json:"id"`
	Owner       string     `json:"owner"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsResolved  bool       `json:"is_resolved"`
	IsArchived  bool       `json:"is_archived"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}
