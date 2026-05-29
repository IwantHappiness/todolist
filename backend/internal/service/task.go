package service

import (
	"context"

	taskdomain "github.com/IwantHappiness/todolist/internal/models"
)

type Repository interface {
	GetAll(ctx context.Context) ([]taskdomain.Task, error)
	GetByID(ctx context.Context, id int) (taskdomain.Task, error)
	GetByCompleted(ctx context.Context, completed bool) ([]taskdomain.Task, error)
	Create(ctx context.Context, task taskdomain.TaskDTO) (taskdomain.Task, error)
	Update(ctx context.Context, id int, task taskdomain.TaskDTO) (taskdomain.Task, error)
	Complete(ctx context.Context, id int, completed taskdomain.CompleteTaskDTO) (taskdomain.Task, error)
	Delete(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllTasks(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetTaskByID(ctx context.Context, id int) (taskdomain.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) GetByCompleted(ctx context.Context, completed bool) ([]taskdomain.Task, error) {
	return s.repo.GetByCompleted(ctx, completed)
}

func (s *Service) CreateTask(ctx context.Context, task taskdomain.TaskDTO) (taskdomain.Task, error) {
	return s.repo.Create(ctx, task)
}

func (s *Service) DeleteTaskByID(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) DeleteAllTask(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}

func (s *Service) UpdateTask(ctx context.Context, id int, task taskdomain.TaskDTO) (taskdomain.Task, error) {
	return s.repo.Update(ctx, id, task)
}

func (s *Service) CompleteTaskStatus(ctx context.Context, id int, completed taskdomain.CompleteTaskDTO) (taskdomain.Task, error) {
	return s.repo.Complete(ctx, id, completed)
}
