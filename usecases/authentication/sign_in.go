package authentication

import (
	"errors"
	"fmt"
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
	cognitoService          services.Cognito
}

func NewSignIn(
	usersRepository repositories.Users,
	sessionsRepository repositories.Sessions,
	tokensService services.Tokens,
	passwordCompareFunction utils.PasswordCompareFn,
	cognitoService services.Cognito,
) SignIn {
	return &signInImpl{
		usersRepository:         usersRepository,
		sessionsRepository:      sessionsRepository,
		tokensService:           tokensService,
		passwordCompareFunction: passwordCompareFunction,
		cognitoService:          cognitoService,
	}
}

func (uc *signInImpl) Execute(email, password string) (string, error) {
	userResp, err := uc.cognitoService.AdminGetUser(email)
	if err == nil {
		// User exists, check verification status
		if userResp.UserStatus == "UNCONFIRMED" {
			// User is unverified, resend verification code
			err = uc.cognitoService.ResendConfirmationCode(email)
			if err != nil {
				return "", fmt.Errorf("failed to resend verification code: %w", err)
			}
			return "", errors.New("user is not verified, verification code has been resent")
		}
	}

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
