package service

import (
	"context"
	"fmt"

	"github.com/BBloggsbott/task-api/internal/models"
	"github.com/BBloggsbott/task-api/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, req *models.CreateTaskRequest) (*models.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	task, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, err
}

func (s *TaskService) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID")
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) ListTasks(ctx context.Context, status string, limit, offset int) ([]*models.Task, error) {
	if limit <= 0 {
		limit = 10 // Default
	}
	if limit > 100 {
		limit = 100 // Max
	}
	if offset < 0 {
		offset = 0
	}

	if status != "" {
		taskStatus := models.TaskStatus(status)
		if !taskStatus.IsValid() {
			return nil, fmt.Errorf("invalid status: %s", status)
		}
	}

	tasks, err := s.repo.GetAll(ctx, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id int64, req *models.UpdateTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID")
	}

	if req.Status != nil && !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", *req.Status)
	}

	task, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid task ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
