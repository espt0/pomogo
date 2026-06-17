package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Task struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	EstimatedPomodoros int64     `json:"estimated_pomodoros"`
	CompletedPomodoros int64     `json:"completed_pomodoros"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
