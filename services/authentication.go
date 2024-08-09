package services

import (
	"errors"
	"relif/bff/entities"
	"relif/bff/utils"
)

type Authentication interface {
	SignUp(data entities.User) (entities.Session, error)
	OrganizationSignUp(data entities.User) (entities.Session, error)
	SignIn(email, password string) (entities.Session, error)
	SignOut(sessionID string) error
	AuthenticateSession(sessionID string) (entities.User, error)
}

type authenticationImpl struct {
	usersService         Users
	sessionsService      Sessions
	organizationsService Organizations
	passwordHashFn       utils.PasswordHashFn
	passwordCompareFn    utils.PasswordCompareFn
}

func NewAuth(
	usersService Users,
	sessionsService Sessions,
	organizationsService Organizations,
	passwordHashFn utils.PasswordHashFn,
	passwordCompareFn utils.PasswordCompareFn,
) Authentication {
	return &authenticationImpl{
		usersService:         usersService,
		sessionsService:      sessionsService,
		organizationsService: organizationsService,
		passwordHashFn:       passwordHashFn,
		passwordCompareFn:    passwordCompareFn,
	}
}

func (service *authenticationImpl) SignUp(data entities.User) (entities.Session, error) {
	hashed, err := service.passwordHashFn(data.Password)

	if err != nil {
		return entities.Session{}, err
	}

	data.Password = hashed
	data.PlatformRole = utils.NoOrgPlatformRole

	user, err := service.usersService.Create(data)

	if err != nil {
		return entities.Session{}, err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authenticationImpl) OrganizationSignUp(data entities.User) (entities.Session, error) {
	hashed, err := service.passwordHashFn(data.Password)

	if err != nil {
		return entities.Session{}, err
	}

	data.Password = hashed
	data.PlatformRole = utils.OrgMemberPlatformRole

	user, err := service.usersService.Create(data)

	if err != nil {
		return entities.Session{}, err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authenticationImpl) SignIn(email, password string) (entities.Session, error) {
	user, err := service.usersService.FindOneByEmail(email)

	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			return entities.Session{}, utils.ErrInvalidCredentials
		}

		return entities.Session{}, err
	}

	if user.OrganizationID != "" {
		organization, err := service.organizationsService.FindOneByID(user.OrganizationID)

		if err != nil {
			return entities.Session{}, err
		}

		if organization.Status == utils.InactiveStatus {
			return entities.Session{}, utils.ErrMemberOfInactiveOrganization
		}
	}

	if err = service.passwordCompareFn(password, user.Password); err != nil {
		return entities.Session{}, utils.ErrInvalidCredentials
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authenticationImpl) SignOut(sessionID string) error {
	return service.sessionsService.DeleteOneBySessionID(sessionID)
}

func (service *authenticationImpl) AuthenticateSession(sessionID string) (entities.User, error) {
	session, err := service.sessionsService.FindOneBySessionID(sessionID)

	if err != nil {
		return entities.User{}, err
	}

	user, err := service.usersService.FindOneByID(session.UserID)

	if err != nil {
		return entities.User{}, err
	}

	if user.OrganizationID != "" {
		organization, err := service.organizationsService.FindOneByID(user.OrganizationID)

		if err != nil {
			return entities.User{}, err
		}

		if organization.Status == utils.InactiveStatus {
			return entities.User{}, utils.ErrMemberOfInactiveOrganization
		}
	}

	return user, nil
}
