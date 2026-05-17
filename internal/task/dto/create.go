package dto

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description"`
	Status      string `json:"status" binding:"required,oneof=todo doing done"`
	Assignee    string `json:"assignee"`
}
