package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type Sessions interface {
	Generate(userId string) (entities.Session, error)
	FindOneById(id string) (entities.Session, error)
	DeleteOneById(id string) error
}

type sessionsImpl struct {
	repository repositories.Sessions
}

func NewSessions(repository repositories.Sessions) Sessions {
	return &sessionsImpl{
		repository: repository,
	}
}

func (service *sessionsImpl) Generate(userId string) (entities.Session, error) {
	session := entities.Session{
		UserID:    userId,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	sessionID, err := service.repository.Generate(session)

	if err != nil {
		return entities.Session{}, err
	}

	session.ID = sessionID

	return session, nil
}

func (service *sessionsImpl) FindOneById(id string) (entities.Session, error) {
	return service.repository.FindOneById(id)
}

func (service *sessionsImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
