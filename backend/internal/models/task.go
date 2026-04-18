package models

import (
	"time"
)

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Completed   bool       `json:"completed"`
	Created     time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Description string     `json:"description"`
}
