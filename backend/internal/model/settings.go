package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type Settings struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	WorkDuration       int64     `json:"work_duration"`
	ShortBreakDuration int64     `json:"short_break_duration"`
	LongBreakDuration  int64     `json:"long_break_duration"`
	LongBreakInterval  int64     `json:"long_break_interval"`
	AutoStartWork      bool      `json:"auto_start_work"`
	AutoStartBreak     bool      `json:"auto_start_break"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
