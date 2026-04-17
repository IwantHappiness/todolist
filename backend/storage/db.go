package storage

import (
	"context"

	"github.com/IwantHappiness/todolist/task"
	"github.com/jackc/pgx/v5"
)

func CreateConnection(ctx context.Context, dbUrl string) (*pgx.Conn, error) {
	return pgx.Connect(ctx, dbUrl)
}

func GetAllTasks(ctx context.Context, conn *pgx.Conn) ([]task.Task, error) {
	query := `
	SELECT * FROM tasks
	ORDER BY id ASC`

	rows, err := conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []task.Task
	for rows.Next() {
		var task task.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.CompletedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
