package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type Sessions interface {
	Generate(userID string) (entities.Session, error)
	FindOneByIDAndUserID(sessionID, userID string) (entities.Session, error)
	DeleteOneByID(sessionID string) error
}

type sessionsImpl struct {
	repository repositories.Sessions
}

func NewSessions(repository repositories.Sessions) Sessions {
	return &sessionsImpl{
		repository: repository,
	}
}

func (service *sessionsImpl) Generate(userID string) (entities.Session, error) {
	session, err := service.repository.FindOneByUserID(userID)

	if err != nil {
		return entities.Session{}, err
	}

	if session.ID == "" {
		data := entities.Session{
			UserID: userID,
		}

		session, err = service.repository.Generate(data)

		if err != nil {
			return entities.Session{}, err
		}
	}

	return session, nil
}

func (service *sessionsImpl) FindOneByIDAndUserID(sessionID, userID string) (entities.Session, error) {
	return service.repository.FindOneByIDAndUserID(sessionID, userID)
}

func (service *sessionsImpl) DeleteOneByID(sessionID string) error {
	return service.repository.DeleteOneByID(sessionID)
}
