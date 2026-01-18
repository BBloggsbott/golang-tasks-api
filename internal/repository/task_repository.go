package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BBloggsbott/task-api/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.CreateTaskRequest) (*models.Task, error) {
	query := `
		INSERT INTO tasks (title, description, status, priority)
		values (?, ?, ?, ?)
	`

	var description sql.NullString
	if task.Description != nil {
		description = sql.NullString{
			String: *task.Description,
			Valid:  true,
		}
	}

	status := task.Status
	if status == "" {
		status = models.TaskStatusPending
	}

	result, err := r.db.ExecContext(ctx, query, task.Title, description, status, task.Priority)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last return id: %w", err)
	}

	// Adds additional query for every write. This has some issues. Will be fixed later
	return r.GetByID(ctx, id)
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	query := `
        SELECT id, title, description, status, priority, created_at, updated_at
        FROM tasks
        WHERE id = ?
    `

	task := &models.Task{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) GetAll(ctx context.Context, status string, limit, offset int) ([]*models.Task, error) {
	query := `
        SELECT id, title, description, status, priority, created_at, updated_at
        FROM tasks
    `
	args := []interface{}{}

	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	tasks := []*models.Task{}
	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int64, req *models.UpdateTaskRequest) (*models.Task, error) {
	// First, check if task exists
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	query := "UPDATE tasks SET "
	args := []interface{}{}
	updates := []string{}

	if req.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *req.Title)
	}

	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, sql.NullString{
			String: *req.Description,
			Valid:  true,
		})
	}

	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", *req.Status)
		}
		updates = append(updates, "status = ?")
		args = append(args, *req.Status)
	}

	if req.Priority != nil {
		updates = append(updates, "priority = ?")
		args = append(args, *req.Priority)
	}

	if len(updates) == 0 {
		return existing, nil
	}

	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"
	args = append(args, id)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Adds additional query for every write. This has some issues. Will be fixed later
	return r.GetByID(ctx, id)
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM tasks WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}
