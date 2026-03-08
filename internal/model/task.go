package model

import (
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the lifecycle state of a task.
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

// Task is the core domain entity.
type Task struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

// CreateTaskRequest is the payload for POST /api/tasks.
type CreateTaskRequest struct {
	Title string `json:"title"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

// UpdateTaskRequest is the payload for PATCH /api/tasks/{id}.
type UpdateTaskRequest struct {
	Status TaskStatus `json:"status"`
}
