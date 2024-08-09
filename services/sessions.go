package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Sessions interface {
	Generate(userID string) (entities.Session, error)
	FindOneBySessionID(sessionID string) (entities.Session, error)
	DeleteOneBySessionID(sessionID string) error
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

func (service *sessionsImpl) Generate(userID string) (entities.Session, error) {
	session := entities.Session{
		UserID:    userID,
		SessionID: service.uuidGenerator(),
	}

	session, err := service.repository.Generate(session)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *sessionsImpl) FindOneBySessionID(sessionID string) (entities.Session, error) {
	return service.repository.FindOneBySessionID(sessionID)
}

func (service *sessionsImpl) DeleteOneBySessionID(sessionID string) error {
	return service.repository.DeleteOneBySessionID(sessionID)
}
