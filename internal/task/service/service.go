package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/imirazimi/graph/internal/infra/metric"
	"github.com/imirazimi/graph/internal/task/entity"
	"github.com/imirazimi/graph/internal/task/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

type TaskService interface {
	Create(ctx context.Context, task *entity.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error)
	List(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, int64, error)
	Update(ctx context.Context, task *entity.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) TaskService {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, task *entity.Task) error {
	task.ID = uuid.New()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	err := s.repo.Create(ctx, task)
	if err != nil {
		return err
	}
	metric.TasksCount.Inc()
	return nil
}
func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

func (s *service) List(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	tasks, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	metric.TasksCount.Set(float64(total))

	return tasks, total, nil
}

func (s *service) Update(ctx context.Context, task *entity.Task) error {
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
	err = s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	metric.TasksCount.Dec()
	return nil
}
