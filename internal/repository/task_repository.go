package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/observer/app/internal/model"
)

// TaskRepository handles all DB operations for tasks.
type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) List(ctx context.Context) ([]*model.Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, status, created_at FROM tasks ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		t := &model.Task{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	t := &model.Task{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, title, status, created_at FROM tasks WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrTaskNotFound
	}
	return t, err
}

func (r *TaskRepository) Insert(ctx context.Context, t *model.Task) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO tasks (title, status) VALUES ($1, $2)
		 RETURNING id, created_at`,
		t.Title, t.Status,
	).Scan(&t.ID, &t.CreatedAt)
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.TaskStatus) (*model.Task, error) {
	t := &model.Task{}
	err := r.pool.QueryRow(ctx,
		`UPDATE tasks SET status = $2 WHERE id = $1
		 RETURNING id, title, status, created_at`,
		id, status,
	).Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrTaskNotFound
	}
	return t, err
}
