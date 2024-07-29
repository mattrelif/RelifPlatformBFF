package models

import (
	"relif/bff/entities"
	"time"
)

type Session struct {
	UserID    string    `bson:"_id"`
	SessionID string    `bson:"session_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func (session *Session) ToEntity() entities.Session {
	return entities.Session{
		UserID:    session.UserID,
		SessionID: session.SessionID,
		ExpiresAt: session.ExpiresAt,
	}
}
