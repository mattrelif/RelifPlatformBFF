package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type Sessions interface {
	Generate(userId string) (entities.Session, error)
	FindOneBySessionId(sessionId string) (entities.Session, error)
	DeleteOneBySessionId(sessionId string) error
}

type sessionsImpl struct {
	repository    repositories.Sessions
	uuidGenerator utils.UuidGenerator
}

func NewSessions(repository repositories.Sessions, uuidGenerator utils.UuidGenerator) Sessions {
	return &sessionsImpl{
		repository:    repository,
		uuidGenerator: uuidGenerator,
	}
}

func (service *sessionsImpl) Generate(userId string) (entities.Session, error) {
	session := entities.Session{
		UserID:    userId,
		SessionID: service.uuidGenerator(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	if err := service.repository.Generate(session); err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *sessionsImpl) FindOneBySessionId(sessionId string) (entities.Session, error) {
	return service.repository.FindOneBySessionId(sessionId)
}

func (service *sessionsImpl) DeleteOneBySessionId(sessionId string) error {
	return service.repository.DeleteOneBySessionId(sessionId)
}
