package repository

import (
	"context"
	"errors"
	"time"

	"github.com/IwantHappiness/todolist/internal/models"
	"github.com/jackc/pgx/v5"
)

type TaskRepository interface {
	GetAll(ctx context.Context) ([]models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, error)
	GetByCompleted(ctx context.Context, completed bool) ([]models.Task, error)
	Create(ctx context.Context, task models.TaskDTO) (models.Task, error)
	Update(ctx context.Context, id int, task models.TaskDTO) (models.Task, error)
	Complete(ctx context.Context, id int, completed models.CompleteTaskDTO) (models.Task, error)
	Delete(ctx context.Context, id int) error
	DeleteAll(ctx context.Context) error
}

type TaskPgRepository struct {
	db *pgx.Conn
}

func NewTaskPgRepository(db *pgx.Conn) TaskRepository {
	return &TaskPgRepository{
		db: db,
	}
}

func (r *TaskPgRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	query := `
	SELECT * FROM tasks
	ORDER BY id ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.Task, 0, 50)

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskPgRepository) GetByID(ctx context.Context, id int) (models.Task, error) {
	query := `
	SELECT * FROM tasks
	WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var task models.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *TaskPgRepository) GetByCompleted(ctx context.Context, completed bool) ([]models.Task, error) {
	query := `
	SELECT * FROM tasks
	WHERE completed = $1`

	rows, err := r.db.Query(ctx, query, completed)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]models.Task, 0, 25)

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskPgRepository) Create(ctx context.Context, taskDTO models.TaskDTO) (models.Task, error) {
	query := `
	INSERT INTO tasks (title, description, created_at)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, completed, created_at, completed_at`

	row := r.db.QueryRow(ctx, query, taskDTO.Title, taskDTO.Description, time.Now())

	var task models.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *TaskPgRepository) Update(ctx context.Context, id int, taskDTO models.TaskDTO) (models.Task, error) {
	query := `
	UPDATE tasks
	SET title = $1, description = $2
	WHERE id = $3
	RETURNING id, title, description, completed, created_at, completed_at`

	var task models.Task
	err := r.db.QueryRow(ctx, query, taskDTO.Title, taskDTO.Description, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Created,
		&task.CompletedAt)

	if err == pgx.ErrNoRows {
		return models.Task{}, errors.New("Task not found")
	}

	return task, err
}

func (r *TaskPgRepository) Complete(ctx context.Context, id int, completed models.CompleteTaskDTO) (models.Task, error) {
	query := `
	UPDATE tasks
	SET completed = $1, completed_at = $2
	WHERE id = $3
	RETURNING id, title, description, completed, created_at, completed_at`

	var task models.Task
	err := r.db.QueryRow(ctx, query, completed.Completed, time.Now(), id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Created,
		&task.CompletedAt)

	if err == pgx.ErrNoRows {
		return models.Task{}, errors.New("Task not found")
	}

	return task, err
}

func (r *TaskPgRepository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *TaskPgRepository) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM tasks`

	_, err := r.db.Exec(ctx, query)
	return err
}
