package authentication

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	usersUseCase "relif/platform-bff/usecases/users"
	"relif/platform-bff/utils"
)

type SignUp interface {
	Execute(data entities.User, locale string) (string, error)
}

type signUpImpl struct {
	sessionsRepository   repositories.Sessions
	tokensService        services.Tokens
	createUserUseCase    usersUseCase.Create
	passwordHashFunction utils.PasswordHashFn
	cognitoService       services.Cognito
}

func NewSignUp(
	sessionsRepository repositories.Sessions,
	tokensService services.Tokens,
	createUserUseCase usersUseCase.Create,
	passwordHashFunction utils.PasswordHashFn,
	cognitoService services.Cognito,
) SignUp {
	return &signUpImpl{
		sessionsRepository:   sessionsRepository,
		tokensService:        tokensService,
		createUserUseCase:    createUserUseCase,
		passwordHashFunction: passwordHashFunction,
		cognitoService:       cognitoService,
	}
}

func (uc *signUpImpl) Execute(data entities.User, locale string) (string, error) {
	hashed, err := uc.passwordHashFunction(data.Password)

	if err != nil {
		return "", err
	}

	data.Password = hashed
	data.PlatformRole = utils.NoOrgPlatformRole
	data.Status = utils.UnverifedStatus

	err = uc.cognitoService.InitiateEmailOrPhoneVerification(data.Email, data.Password, data.FirstName, locale)
	if err != nil {
		// Handle error (you might want to delete the created user if this fails)
		return "", err
	}

	user, err := uc.createUserUseCase.Execute(data)

	if err != nil {
		return "", err
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
