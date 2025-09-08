package models

import "time"

// CreateTaskReq struct
type CreateTaskReq struct {
	ID          string    `json:"-"`
	UserID      string    `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
}

// CreateTaskResp struct
type CreateTaskResp struct {
	ID string `json:"id"`
}

// GetTaskResp struct
type GetTaskResp struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
