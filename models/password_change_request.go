package models

import (
	"relif/bff/entities"
	"time"
)

type PasswordChangeRequest struct {
	UserID    string    `bson:"_id"`
	Code      string    `bson:"code"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func (request *PasswordChangeRequest) ToEntity() entities.PasswordChangeRequest {
	return entities.PasswordChangeRequest{
		UserID:    request.UserID,
		Code:      request.Code,
		ExpiresAt: request.ExpiresAt,
	}
}
