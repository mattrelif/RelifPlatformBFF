package authentication

import (
	"errors"
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type SignIn interface {
	Execute(email, password string) (string, error)
}

type signInImpl struct {
	usersRepository         repositories.Users
	sessionsRepository      repositories.Sessions
	tokensService           services.Tokens
	passwordCompareFunction utils.PasswordCompareFn
}

func NewSignIn(
	usersRepository repositories.Users,
	sessionsRepository repositories.Sessions,
	tokensService services.Tokens,
	passwordCompareFunction utils.PasswordCompareFn,
) SignIn {
	return &signInImpl{
		usersRepository:         usersRepository,
		sessionsRepository:      sessionsRepository,
		tokensService:           tokensService,
		passwordCompareFunction: passwordCompareFunction,
	}
}

func (uc *signInImpl) Execute(email, password string) (string, error) {
	user, err := uc.usersRepository.FindOneAndLookupByEmail(email)

	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return "", utils.ErrInvalidCredentials
		}

		return "", err
	}

	if err = guards.CanAccessPlatform(user); err != nil {
		return "", err
	}

	if err = uc.passwordCompareFunction(password, user.Password); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	session := entities.Session{
		UserID: user.ID,
	}

	session, err = uc.sessionsRepository.Generate(session)

	if err != nil {
		return "", err
	}

	token, err := uc.tokensService.SignToken(user, session)

	if err != nil {
		return "", err
	}

	return token, nil
}
