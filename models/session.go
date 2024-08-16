package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"time"
)

type Session struct {
	ID        string    `bson:"_id,omitempty"`
	UserID    string    `bson:"user_id,omitempty"`
	ExpiresAt time.Time `bson:"expires_at,omitempty"`
}

func (session *Session) ToEntity() entities.Session {
	return entities.Session{
		ID:        session.ID,
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt,
	}
}

func NewSession(entity entities.Session) Session {
	return Session{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    entity.UserID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
}
