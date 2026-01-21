package models

// Account представляет данные пользователя
type Account struct {
	UserID   string `json:"user_id"`
	Secret   string `json:"secret"`
}
