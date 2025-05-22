package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title" validate:"required,min=1,max=255"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Completed   bool      `json:"completed"`
}

type TaskRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description"`
}

type TaskResponse struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Completed   bool      `json:"completed"`
}
