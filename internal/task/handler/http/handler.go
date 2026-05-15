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

func (h *Handler) RegisterRoutes(router *gin.Engine) {
    
}

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

func (h *Handler) List(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    offset := (page - 1) * limit

    filter := repository.TaskFilter{
        Status:   c.Query("status"),
        Assignee: c.Query("assignee"),
        Limit:    limit,
        Offset:   offset,
    }

    tasks, err := h.service.List(c.Request.Context(), filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": tasks,
        "page": page,
        "limit": limit,
    })
}

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