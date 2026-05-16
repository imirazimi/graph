package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/imirazimi/graph/internal/task/entity"
	"github.com/imirazimi/graph/internal/task/repository"
	"github.com/imirazimi/graph/internal/task/service"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockRepository) List(ctx context.Context, filter repository.TaskFilter) ([]entity.Task, int64, error) {
	args := m.Called(ctx, filter)

	var tasks []entity.Task
	if args.Get(0) != nil {
		tasks = args.Get(0).([]entity.Task)
	}

	var total int64
	if args.Get(1) != nil {
		total = args.Get(1).(int64)
	}

	return tasks, total, args.Error(2)
}

func (m *MockRepository) Update(ctx context.Context, task *entity.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateTask(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	task := &entity.Task{
		Title:  "test",
		Status: "todo",
	}

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	err := svc.Create(context.Background(), task)

	assert.NoError(t, err)
	assert.NotEmpty(t, task.ID)

	repo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	expectedTask := &entity.Task{
		ID:     taskID,
		Title:  "task",
		Status: "todo",
	}

	repo.On("GetByID", mock.Anything, taskID).
		Return(expectedTask, nil)

	result, err := svc.GetByID(context.Background(), taskID)

	assert.NoError(t, err)
	assert.Equal(t, expectedTask.ID, result.ID)

	repo.AssertExpectations(t)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(nil, errors.New("not found"))

	result, err := svc.GetByID(context.Background(), taskID)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)

	repo.AssertExpectations(t)
}

func TestListTasks(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	expectedTasks := []entity.Task{
		{Title: "task1", Status: "todo"},
	}

	filter := repository.TaskFilter{
		Limit: 10,
	}

	repo.On("List", mock.Anything, filter).
		Return(expectedTasks, int64(1), nil)

	result, total, err := svc.List(context.Background(), filter)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), total)

	repo.AssertExpectations(t)
}

func TestListTasks_DefaultLimit(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	repo.On("List", mock.Anything, mock.Anything).
		Return([]entity.Task{}, int64(0), nil)

	_, _, err := svc.List(context.Background(), repository.TaskFilter{})

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	existingTask := &entity.Task{
		ID:     taskID,
		Title:  "old",
		Status: "todo",
	}

	repo.On("GetByID", mock.Anything, taskID).
		Return(existingTask, nil)

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil)

	err := svc.Update(context.Background(), &entity.Task{
		ID:     taskID,
		Title:  "new",
		Status: "doing",
	})

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUpdateTask_NotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(nil, errors.New("not found"))

	err := svc.Update(context.Background(), &entity.Task{
		ID: taskID,
	})

	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)

	repo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	existingTask := &entity.Task{ID: taskID}

	repo.On("GetByID", mock.Anything, taskID).
		Return(existingTask, nil)

	repo.On("Delete", mock.Anything, taskID).
		Return(nil)

	err := svc.Delete(context.Background(), taskID)

	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestDeleteTask_NotFound(t *testing.T) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(nil, errors.New("not found"))

	err := svc.Delete(context.Background(), taskID)

	assert.Error(t, err)
	assert.Equal(t, service.ErrTaskNotFound, err)

	repo.AssertExpectations(t)
}

func BenchmarkCreateTask(b *testing.B) {
	repo := new(MockRepository)
	svc := service.NewService(repo)

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	for i := 0; i < b.N; i++ {
		_ = svc.Create(context.Background(), &entity.Task{
			Title:  "benchmark",
			Status: "todo",
		})
	}
}