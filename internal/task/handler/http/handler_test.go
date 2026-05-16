package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
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

/* ---------------- tests ---------------- */

func TestCreateTask(t *testing.T) {
	router, repo := setupHandler()

	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	payload := dto.CreateTaskRequest{
		Title:       "learn golang",
		Description: "practice",
		Status:      "todo",
		Assignee:    "amir",
	}

	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	repo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(&entity.Task{ID: taskID}, nil)

	req := httptest.NewRequest(http.MethodGet, "/tasks/"+taskID.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}

func TestListTasks(t *testing.T) {
	router, repo := setupHandler()

	repo.On("List", mock.Anything, mock.Anything).
		Return([]entity.Task{
			{
				Title:    "task1",
				Status:   "todo",
				Assignee: "amir",
			},
		}, int64(1), nil)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks?status=todo&assignee=amir&page=1&limit=10",
		nil,
	)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, float64(10), resp["limit"])
	assert.Equal(t, float64(1), resp["total"])

	data := resp["data"].([]interface{})
	assert.Len(t, data, 1)

	repo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	router, repo := setupHandler()

	taskID := uuid.New()

	repo.On("GetByID", mock.Anything, taskID).
		Return(&entity.Task{ID: taskID}, nil)

	repo.On("Delete", mock.Anything, taskID).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+taskID.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	repo.AssertExpectations(t)
}