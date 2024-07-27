package models

import (
	"relif/bff/entities"
	"time"
)

type Session struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	ExpiresAt time.Time `bson:"expires_at"`
}

func (session *Session) ToEntity() entities.Session {
	return entities.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt,
	}
}
