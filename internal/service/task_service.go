package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/observer/app/internal/model"
	"github.com/observer/app/internal/repository"
)

// TaskService implements business logic for tasks.
type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) List(ctx context.Context) ([]*model.Task, error) {
	return s.repo.List(ctx)
}

func (s *TaskService) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) Create(ctx context.Context, req *model.CreateTaskRequest) (*model.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	task := &model.Task{
		Title:  req.Title,
		Status: model.TaskStatusPending,
	}
	if err := s.repo.Insert(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, id uuid.UUID, req *model.UpdateTaskRequest) (*model.Task, error) {
	switch req.Status {
	case model.TaskStatusPending, model.TaskStatusInProgress, model.TaskStatusDone:
		// valid
	default:
		return nil, model.ErrInvalidStatus
	}
	return s.repo.UpdateStatus(ctx, id, req.Status)
}
