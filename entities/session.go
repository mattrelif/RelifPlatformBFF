package entities

import (
	"time"
)

type Session struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}
