package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/imirazimi/graph/internal/task/dto"
	"github.com/imirazimi/graph/internal/task/entity"
	handler "github.com/imirazimi/graph/internal/task/handler/http"
	"github.com/imirazimi/graph/internal/task/repository"
	"github.com/imirazimi/graph/internal/task/service"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(
	ctx context.Context,
	task *entity.Task,
) error {
	args := m.Called(ctx, task)

	return args.Error(0)
}

func (m *MockRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*entity.Task, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Task), args.Error(1)
}

func (m *MockRepository) List(
	ctx context.Context,
	filter repository.TaskFilter,
) ([]entity.Task, error) {
	args := m.Called(ctx, filter)

	return args.Get(0).([]entity.Task), args.Error(1)
}

func (m *MockRepository) Update(
	ctx context.Context,
	task *entity.Task,
) error {
	args := m.Called(ctx, task)

	return args.Error(0)
}

func (m *MockRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func setupHandler() (*gin.Engine, *MockRepository) {
	gin.SetMode(gin.TestMode)

	repo := new(MockRepository)

	svc := service.NewService(repo)

	h := handler.NewHandler(svc)

	router := gin.Default()

	router.POST("/tasks", h.Create)
	router.GET("/tasks/:id", h.GetByID)
	router.GET("/tasks", h.List)
	router.PUT("/tasks/:id", h.Update)
	router.DELETE("/tasks/:id", h.Delete)

	return router, repo
}

func TestCreateTask(t *testing.T) {
	router, repo := setupHandler()

	repo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	payload := dto.CreateTaskRequest{
		Title:       "learn golang",
		Description: "practice",
		Status:      "todo",
		Assignee:    "amir",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	repo.AssertExpectations(t)
}

func TestCreateTask_InvalidBody(t *testing.T) {
	router, _ := setupHandler()

	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		bytes.NewBuffer([]byte(`invalid-json`)),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetTaskByID(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(&entity.Task{
			ID:     taskID,
			Title:  "task",
			Status: "todo",
		}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+taskID.String(),
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	repo.AssertExpectations(t)
}

func TestGetTaskByID_InvalidUUID(t *testing.T) {
	router, _ := setupHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/invalid",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetTaskByID_NotFound(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(nil, errors.New("not found"))

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/"+taskID.String(),
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	repo.AssertExpectations(t)
}

func TestListTasks(t *testing.T) {
	router, repo := setupHandler()

	repo.On("List", mock.Anything, mock.Anything).
		Return([]entity.Task{
			{
				Title: "task1",
			},
		}, nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	repo.AssertExpectations(t)
}

func TestUpdateTask(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(&entity.Task{
			ID: taskID,
		}, nil)

	repo.On("Update", mock.Anything, mock.Anything).
		Return(nil)

	payload := dto.UpdateTaskRequest{
		Title: "updated",
		Status: "doing",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/"+taskID.String(),
		bytes.NewBuffer(body),
	)

	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	repo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(&entity.Task{
			ID: taskID,
		}, nil)

	repo.On("Delete", mock.Anything, taskID).
		Return(nil)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/"+taskID.String(),
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	repo.AssertExpectations(t)
}