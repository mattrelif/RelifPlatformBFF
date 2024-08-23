package authentication

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type SignOut interface {
	Execute(session entities.Session) error
}

type signOutImpl struct {
	sessionsRepository repositories.Sessions
}

func NewSignOut(sessionsRepository repositories.Sessions) SignOut {
	return &signOutImpl{sessionsRepository: sessionsRepository}
}

func (uc *signOutImpl) Execute(session entities.Session) error {
	return uc.sessionsRepository.DeleteOneByID(session.ID)
}
