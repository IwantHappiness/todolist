package task

import (
	"errors"
	"time"
)

type CompleteTaskDTO struct {
	Completed bool `json:"completed"`
}

type TaskDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (t TaskDTO) Validate() error {
	if t.Title == "" {
		return errors.New("Title is empty")
	}

	return nil
}

type ErrorDTO struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
