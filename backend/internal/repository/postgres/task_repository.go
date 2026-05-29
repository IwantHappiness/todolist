package repository

import (
	"context"
	"errors"
	"time"

	taskdomain "github.com/IwantHappiness/todolist/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) GetAll(ctx context.Context) ([]taskdomain.Task, error) {
	query := `
	SELECT * FROM tasks
	ORDER BY id ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]taskdomain.Task, 0, 50)

	for rows.Next() {
		var task taskdomain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (taskdomain.Task, error) {
	query := `
	SELECT * FROM tasks
	WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)

	var task taskdomain.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
		return taskdomain.Task{}, err
	}

	return task, nil
}

func (r *Repository) GetByCompleted(ctx context.Context, completed bool) ([]taskdomain.Task, error) {
	query := `
	SELECT * FROM tasks
	WHERE completed = $1`

	rows, err := r.pool.Query(ctx, query, completed)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tasks := make([]taskdomain.Task, 0, 25)

	for rows.Next() {
		var task taskdomain.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *Repository) Create(ctx context.Context, taskDTO taskdomain.TaskDTO) (taskdomain.Task, error) {
	query := `
	INSERT INTO tasks (title, description, created_at)
	VALUES ($1, $2, $3)
	RETURNING id, title, description, completed, created_at, completed_at`

	row := r.pool.QueryRow(ctx, query, taskDTO.Title, taskDTO.Description, time.Now())

	var task taskdomain.Task
	if err := row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.Created, &task.CompletedAt); err != nil {
		return taskdomain.Task{}, err
	}

	return task, nil
}

func (r *Repository) Update(ctx context.Context, id int, taskDTO taskdomain.TaskDTO) (taskdomain.Task, error) {
	query := `
	UPDATE tasks
	SET title = $1, description = $2
	WHERE id = $3
	RETURNING id, title, description, completed, created_at, completed_at`

	var task taskdomain.Task
	err := r.pool.QueryRow(ctx, query, taskDTO.Title, taskDTO.Description, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Created,
		&task.CompletedAt)

	if err == pgx.ErrNoRows {
		return taskdomain.Task{}, errors.New("task not found")
	}

	return task, err
}

func (r *Repository) Complete(ctx context.Context, id int, completed taskdomain.CompleteTaskDTO) (taskdomain.Task, error) {
	query := `
	UPDATE tasks
	SET completed = $1, completed_at = $2
	WHERE id = $3
	RETURNING id, title, description, completed, created_at, completed_at`

	var task taskdomain.Task
	err := r.pool.QueryRow(ctx, query, completed.Completed, time.Now(), id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.Created,
		&task.CompletedAt)

	if err == pgx.ErrNoRows {
		return taskdomain.Task{}, errors.New("task not found")
	}

	return task, err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM tasks
	WHERE id = $1`

	_, err := r.pool.Exec(ctx, query, id)
	return err
}

func (r *Repository) DeleteAll(ctx context.Context) error {
	query := `DELETE FROM tasks`

	_, err := r.pool.Exec(ctx, query)
	return err
}
