package repository

import (
    "context"

    "github.com/google/uuid"
    "github.com/imirazimi/graph/internal/task/entity"
)

type TaskRepository interface {
    Create(ctx context.Context, task *entity.Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
    List(ctx context.Context, filter TaskFilter) ([]entity.Task, error)
    Update(ctx context.Context, task *entity.Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type TaskFilter struct {
    Status   string
    Assignee string
    Limit    int
    Offset   int
}