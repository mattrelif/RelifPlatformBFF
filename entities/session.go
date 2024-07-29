package entities

import (
	"time"
)

type Session struct {
	UserID    string
	SessionID string
	ExpiresAt time.Time
}
