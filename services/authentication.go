package services

import (
	"errors"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

type Authentication interface {
	SignUp(data entities.User) (string, error)
	OrganizationSignUp(data entities.User) (string, error)
	SignIn(email, password string) (string, error)
	SignOut(sessionID string) error
	AuthenticateToken(token string) (entities.User, entities.Session, error)
}

type authenticationImpl struct {
	usersService         Users
	sessionsService      Sessions
	organizationsService Organizations
	tokensService        Tokens
	passwordHashFn       utils.PasswordHashFn
	passwordCompareFn    utils.PasswordCompareFn
}

func NewAuthentication(
	usersService Users,
	sessionsService Sessions,
	organizationsService Organizations,
	tokensService Tokens,
	passwordHashFn utils.PasswordHashFn,
	passwordCompareFn utils.PasswordCompareFn,
) Authentication {
	return &authenticationImpl{
		usersService:         usersService,
		sessionsService:      sessionsService,
		organizationsService: organizationsService,
		tokensService:        tokensService,
		passwordHashFn:       passwordHashFn,
		passwordCompareFn:    passwordCompareFn,
	}
}

func (service *authenticationImpl) SignUp(data entities.User) (string, error) {
	hashed, err := service.passwordHashFn(data.Password)

	if err != nil {
		return "", err
	}

	data.Password = hashed
	data.PlatformRole = utils.NoOrgPlatformRole

	user, err := service.usersService.Create(data)

	if err != nil {
		return "", err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return "", err
	}

	token, err := service.tokensService.SignToken(user, session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *authenticationImpl) OrganizationSignUp(data entities.User) (string, error) {
	hashed, err := service.passwordHashFn(data.Password)

	if err != nil {
		return "", err
	}

	data.Password = hashed
	data.PlatformRole = utils.OrgMemberPlatformRole

	user, err := service.usersService.Create(data)

	if err != nil {
		return "", err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return "", err
	}

	token, err := service.tokensService.SignToken(user, session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *authenticationImpl) SignIn(email, password string) (string, error) {
	user, err := service.usersService.FindOneByEmail(email)

	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return "", utils.ErrInvalidCredentials
		}

		return "", err
	}

	if user.OrganizationID != "" {
		var organization entities.Organization

		organization, err = service.organizationsService.FindOneByID(user.OrganizationID)

		if err != nil {
			return "", err
		}

		if organization.Status == utils.InactiveStatus {
			return "", utils.ErrMemberOfInactiveOrganization
		}
	}

	if err = service.passwordCompareFn(password, user.Password); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return "", err
	}

	token, err := service.tokensService.SignToken(user, session)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *authenticationImpl) SignOut(sessionID string) error {
	return service.sessionsService.DeleteOneByID(sessionID)
}

func (service *authenticationImpl) AuthenticateToken(token string) (entities.User, entities.Session, error) {
	sessionID, userID, err := service.tokensService.ParseToken(token)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	session, err := service.sessionsService.FindOneByIDAndUserID(sessionID, userID)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	user, err := service.usersService.FindOneByID(session.UserID)

	if err != nil {
		return entities.User{}, entities.Session{}, err
	}

	if user.OrganizationID != "" {
		var organization entities.Organization

		organization, err = service.organizationsService.FindOneByID(user.OrganizationID)

		if err != nil {
			return entities.User{}, entities.Session{}, err
		}

		if organization.Status == utils.InactiveStatus {
			return entities.User{}, entities.Session{}, utils.ErrMemberOfInactiveOrganization
		}
	}

	return user, session, nil
}
