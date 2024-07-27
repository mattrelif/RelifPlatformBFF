package entities

import "time"

type PasswordChangeRequest struct {
	ID        string
	UserID    string
	ExpiresAt time.Time
}
