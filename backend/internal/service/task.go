package service

import (
	"context"

	"github.com/IwantHappiness/todolist/internal/models"
	"github.com/IwantHappiness/todolist/internal/repository"
)

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *taskService) GetTaskByID(ctx context.Context, id int) (models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) GetByCompleted(ctx context.Context, completed bool) ([]models.Task, error) {
	return s.repo.GetByCompleted(ctx, completed)
}

func (s *taskService) CreateTask(ctx context.Context, task models.TaskDTO) (models.Task, error) {
	return s.repo.Create(ctx, task)
}

func (s *taskService) DeleteTaskByID(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskService) DeleteAllTask(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}

func (s *taskService) UpdateTask(ctx context.Context, id int, task models.TaskDTO) (models.Task, error) {
	return s.repo.Update(ctx, id, task)
}

func (s *taskService) CompleteTaskStatus(ctx context.Context, id int, completed models.CompleteTaskDTO) (models.Task, error) {
	return s.repo.Complete(ctx, id, completed)
}
