package models

import "time"

// Notification user Reminder notification
type Notification struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TaskID    int64     `json:"task_id"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
