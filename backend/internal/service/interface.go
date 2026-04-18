package service

import (
	"context"

	"github.com/IwantHappiness/todolist/internal/models"
)

type TaskService interface {
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetTaskByID(ctx context.Context, id int) (models.Task, error)
	GetByCompleted(ctx context.Context, completed bool) ([]models.Task, error)
	CreateTask(ctx context.Context, task models.TaskDTO) (models.Task, error)
	DeleteTaskByID(ctx context.Context, id int) error
	DeleteAllTask(ctx context.Context) error
	UpdateTask(ctx context.Context, id int, task models.TaskDTO) (models.Task, error)
	CompleteTaskStatus(ctx context.Context, id int, completed models.CompleteTaskDTO) (models.Task, error)
}
