package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Task struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	UserID             uuid.UUID `json:"user_id" db:"user_id "`
	Title              string    `json:"title" db:"title"`
	Description        string    `json:"description" db:"description"`
	Status             string    `json:"status" db:"status"`
	Active             bool      `json:"active" db:"active"`
	EstimatedPomodoros int64     `json:"estimated_pomodoros" db:"estimated_pomodoros"`
	CompletedPomodoros int64     `json:"completed_pomodoros" db:"completed_pomodoros"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTaskRequest struct {
	Title              string `json:"title" validate:"required,min=3,max=200"`
	Description        string `json:"description" validate:"required,max=1000"`
	EstimatedPomodoros int64  `json:"estimated_pomodoros" validate:"required,min=1"`
}

type UpdateTaskRequest struct {
	Title              *string `json:"title" validate:"omitempty,min=3,max=200"`
	Description        *string `json:"description" validate:"omitempty,max=1000"`
	Status             *string `json:"status" validate:"omitempty,oneof=pending in_progress completed archived"`
	EstimatedPomodoros *int64  `json:"estimated_pomodoros" validate:"omitempty,min=1"`
	CompletedPomodoros *int64  `json:"completed_pomodoros" validate:"omitempty,min=1"`
}
