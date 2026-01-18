package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Status      TaskStatus     `json:"status"`
	Priority    int            `json:"priority"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	Status      TaskStatus `json:"status,omitempty"`
	Priority    int        `json:"priority,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string     `json:"title,omitempty" binding:"omitempty,min=1,max=255"`
	Description *string     `json:"description,omitempty"`
	Status      *TaskStatus `json:"status,omitempty"`
	Priority    *int        `json:"priority,omitempty"`
}

func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusCompleted, TaskStatusInProgress, TaskStatusPending:
		return true
	}
	return false
}

func (s TaskStatus) String() string {
	return string(s)
}

func (t Task) MarshalJSON() ([]byte, error) {
	type Alias Task

	var description *string
	if t.Description.Valid {
		description = &t.Description.String
	}

	return json.Marshal(&struct {
		*Alias
		Description *string `json:"description"`
	}{
		Alias:       (*Alias)(&t),
		Description: description,
	})
}

func (t *Task) UnMarshalJSON(data []byte) error {
	type Alias Task

	aux := &struct {
		Description *string `json:"description"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Description != nil {
		t.Description = sql.NullString{
			String: *aux.Description,
			Valid:  true,
		}
	} else {
		t.Description = sql.NullString{Valid: false}
	}

	return nil

}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(r.Title) > 255 {
		return fmt.Errorf("title should be less than 255 characters")
	}

	if r.Status == "" {
		r.Status = TaskStatusPending
	}

	if !r.Status.IsValid() {
		return fmt.Errorf("invalid status for task: %s", r.Status)
	}

	return nil
}
