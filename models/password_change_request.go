package models

import (
	"relif/bff/entities"
	"time"
)

type PasswordChangeRequest struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func (request *PasswordChangeRequest) ToEntity() entities.PasswordChangeRequest {
	return entities.PasswordChangeRequest{
		ID:        request.ID,
		UserID:    request.UserID,
		ExpiresAt: request.ExpiresAt,
	}
}
