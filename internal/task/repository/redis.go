package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imirazimi/graph/internal/task/entity"
	"github.com/imirazimi/graph/internal/infra/redis"
)

type cacheRepository struct {
	redis redis.RedisClient
	next  TaskRepository
}

func NewCacheRepository(redis redis.RedisClient, next TaskRepository) TaskRepository {
	return &cacheRepository{
		redis: redis,
		next:  next,
	}
}

type taskListCache struct {
	Data  []entity.Task `json:"data"`
	Total int64         `json:"total"`
}

func normalize(v string) string {
	if v == "" {
		return "all"
	}
	return v
}

func generateTaskListKey(filter TaskFilter) string {
	return fmt.Sprintf(
		"tasks:list:status=%s:assignee=%s:limit=%d:offset=%d",
		normalize(filter.Status),
		normalize(filter.Assignee),
		filter.Limit,
		filter.Offset,
	)
}

// ---------------- LIST ----------------

func (c *cacheRepository) List(ctx context.Context, filter TaskFilter) ([]entity.Task, int64, error) {
	cacheKey := generateTaskListKey(filter)

	// cache hit
	cachedValue, err := c.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var cached taskListCache
		if json.Unmarshal([]byte(cachedValue), &cached) == nil {
			return cached.Data, cached.Total, nil
		}
	}

	// DB fallback
	tasks, total, err := c.next.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// set cache
	payload := taskListCache{
		Data:  tasks,
		Total: total,
	}

	if data, err := json.Marshal(payload); err == nil {
		_ = c.redis.Set(ctx, cacheKey, data, 5*time.Minute).Err()
	}

	return tasks, total, nil
}

// ---------------- CRUD passthrough ----------------

func (c *cacheRepository) Create(ctx context.Context, task *entity.Task) error {
	if err := c.next.Create(ctx, task); err != nil {
		return err
	}
	return nil
}

func (c *cacheRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Task, error) {
	return c.next.GetByID(ctx, id)
}

func (c *cacheRepository) Update(ctx context.Context, task *entity.Task) error {
	return c.next.Update(ctx, task)
}

func (c *cacheRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return c.next.Delete(ctx, id)
}