package models

import (
	"relif/bff/entities"
	"time"
)

type Session struct {
	UserID    string    `bson:"_id,omitempty"`
	SessionID string    `bson:"session_id,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

func (session *Session) ToEntity() entities.Session {
	return entities.Session{
		UserID:    session.UserID,
		SessionID: session.SessionID,
		ExpiresAt: session.ExpiresAt,
	}
}

func NewSession(entity entities.Session) Session {
	return Session{
		UserID:    entity.UserID,
		SessionID: entity.SessionID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
}
