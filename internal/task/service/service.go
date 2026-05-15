package service

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/imirazimi/graph/internal/task/model"
    "github.com/imirazimi/graph/internal/task/repository"
)

var (
    ErrTaskNotFound = errors.New("task not found")
)

type TaskService interface {
    Create(ctx context.Context, task *model.Task) error
    GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
    List(ctx context.Context, filter repository.TaskFilter) ([]model.Task, error)
    Update(ctx context.Context, task *model.Task) error
    Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
    repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, task *model.Task) error {
    task.ID = uuid.New()
    task.CreatedAt = time.Now()
    task.UpdatedAt = time.Now()

    return s.repo.Create(ctx, task)
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
    task, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, ErrTaskNotFound
    }

    return task, nil
}

func (s *service) List(ctx context.Context, filter repository.TaskFilter) ([]model.Task, error) {
    if filter.Limit <= 0 {
        filter.Limit = 10
    }

    return s.repo.List(ctx, filter)
}

func (s *service) Update(ctx context.Context, task *model.Task) error {
    existingTask, err := s.repo.GetByID(ctx, task.ID)
    if err != nil {
        return ErrTaskNotFound
    }

    existingTask.Title = task.Title
    existingTask.Description = task.Description
    existingTask.Status = task.Status
    existingTask.Assignee = task.Assignee
    existingTask.UpdatedAt = time.Now()

    return s.repo.Update(ctx, existingTask)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return ErrTaskNotFound
    }

    return s.repo.Delete(ctx, id)
}