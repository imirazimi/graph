package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/imirazimi/graph/internal/task/dto"
	"github.com/imirazimi/graph/internal/task/entity"
	"github.com/imirazimi/graph/internal/task/repository"
	"github.com/imirazimi/graph/internal/task/service"
)

type Handler struct {
	service service.TaskService
}

func NewHandler(service service.TaskService) Handler {
	return Handler{service: service}
}

// Create godoc
// @Summary Create task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task payload"
// @Success 201 {object} entity.Task
// @Failure 400 {object} map[string]interface{}
// @Router /tasks [post]
func (h *Handler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Assignee:    req.Assignee,
	}

	err := h.service.Create(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetByID godoc
// @Summary Get task
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} entity.Task
// @Failure 404 {object} map[string]interface{}
// @Router /tasks/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// List godoc
// @Summary      List tasks
// @Description  Get paginated list of tasks with optional filters
// @Tags         tasks
// @Produce      json
// @Param        status    query     string  false  "Filter by status"
// @Param        assignee  query     string  false  "Filter by assignee"
// @Param        page      query     int     false  "Page number (default: 1)"
// @Param        limit     query     int     false  "Items per page (default: 10)"
// @Success      200       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]string
// @Router       /tasks [get]
func (h *Handler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	filter := repository.TaskFilter{
		Status:   c.Query("status"),
		Assignee: c.Query("assignee"),
		Limit:    limit,
		Offset:   offset,
	}

	tasks, total, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

// Update godoc
// @Summary Update task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body dto.UpdateTaskRequest true "Task payload"
// @Success 200 {object} map[string]interface{}
// @Router /tasks/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	var req dto.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	task := &entity.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Assignee:    req.Assignee,
	}

	err = h.service.Update(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task updated successfully",
	})
}

// Delete godoc
// @Summary Delete task
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]interface{}
// @Router /tasks/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
