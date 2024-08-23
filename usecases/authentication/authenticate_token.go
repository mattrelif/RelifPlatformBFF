package authentication

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
)

type AuthenticateToken interface {
	Execute(token string) (entities.User, entities.Session, error)
}

type authenticateTokenImpl struct {
	usersRepository    repositories.Users
	sessionsRepository repositories.Sessions
	tokensService      services.Tokens
}

func NewAuthenticateToken(
	usersRepository repositories.Users,
	sessionsRepository repositories.Sessions,
	tokensService services.Tokens,
) AuthenticateToken {
	return &authenticateTokenImpl{
		usersRepository:    usersRepository,
		sessionsRepository: sessionsRepository,
		tokensService:      tokensService,
	}
}

func (uc *authenticateTokenImpl) Execute(token string) (entities.User, entities.Session, error) {
	sessionID, userID, err := uc.tokensService.ParseToken(token)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	session, err := uc.sessionsRepository.FindOneByIDAndUserID(sessionID, userID)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	user, err := uc.usersRepository.FindOneAndLookupByID(session.UserID)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	if err = guards.CanAccessPlatform(user); err != nil {
		return entities.User{}, entities.Session{}, err
	}

	return user, session, nil
}
