package entities

import "time"

type PasswordChangeRequest struct {
	UserID    string
	Code      string
	ExpiresAt time.Time
}
