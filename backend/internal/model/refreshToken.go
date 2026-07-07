package model

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Revoked   bool      `json:"revoked" db:"revoked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
