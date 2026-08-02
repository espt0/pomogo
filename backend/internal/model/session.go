package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Session struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	TaskID          uuid.UUID  `json:"task_id" db:"task_id"`
	Type            string     `json:"type" db:"type"`
	DurationMinutes int64      `json:"duration" db:"duration"`
	StartedAt       time.Time  `json:"started_at" db:"started_at"`
	EndedAt         *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	Completed       bool       `json:"completed" db:"completed"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
}

type CreateSessionRequest struct {
	TaskID uuid.UUID `json:"task_id" validate:"required"`
	Type   string    `json:"type" validate:"required,oneof=work short_break long_break"`
}
