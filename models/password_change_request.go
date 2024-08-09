package models

import (
	"relif/platform-bff/entities"
	"time"
)

type PasswordChangeRequest struct {
	UserID    string    `bson:"_id,omitempty"`
	Code      string    `bson:"code,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

func (request *PasswordChangeRequest) ToEntity() entities.PasswordChangeRequest {
	return entities.PasswordChangeRequest{
		UserID:    request.UserID,
		Code:      request.Code,
		ExpiresAt: request.ExpiresAt,
	}
}

func NewPasswordChangeRequest(entity entities.PasswordChangeRequest) PasswordChangeRequest {
	return PasswordChangeRequest{
		Code:      entity.Code,
		ExpiresAt: time.Now().Add(time.Hour * 4),
	}
}
