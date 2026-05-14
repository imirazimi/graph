package dto

import (
    "time"

    "github.com/google/uuid"
)

type TaskResponse struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Status      string    `json:"status"`
    Assignee    string    `json:"assignee"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}