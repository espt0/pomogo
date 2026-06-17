package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TaskID          uuid.UUID `json:"task_id"`
	Type            string    `json:"type"`
	DurationMinutes int64     `json:"duration"`
	StartedAt       time.Time `json:"started_at"`
	EndedAt         time.Time `json:"ended_at"`
	Completed       bool      `json:"completed"`
	CreatedAt       time.Time `json:"created_at"`
}
