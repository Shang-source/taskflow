package models

import "time"

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsDone      bool       `json:"is_done"`
	UserID      int64      `json:"user_id"`
	DueDate     *time.Time `json:"due_date,omitempty"` // Expiration time
	Reminded    bool       `json:"reminded"`           //
}
