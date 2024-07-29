package services

import (
	"relif/bff/entities"
	"relif/bff/utils"
)

type Auth interface {
	SignUp(user entities.User) (entities.Session, error)
	OrganizationSignUp(user entities.User) (entities.Session, error)
	SignIn(email, password string) (entities.Session, error)
	SignOut(sessionId string) error
	AuthenticateSession(sessionId string) (entities.User, error)
}

type authImpl struct {
	usersService      Users
	sessionsService   Sessions
	passwordHashFn    utils.PasswordHashFn
	passwordCompareFn utils.PasswordCompareFn
}

func NewAuth(
	usersService Users,
	sessionsService Sessions,
	passwordHashFn utils.PasswordHashFn,
	passwordCompareFn utils.PasswordCompareFn,
) Auth {
	return &authImpl{
		usersService:      usersService,
		sessionsService:   sessionsService,
		passwordHashFn:    passwordHashFn,
		passwordCompareFn: passwordCompareFn,
	}
}

func (service *authImpl) SignUp(user entities.User) (entities.Session, error) {
	hashed, err := service.passwordHashFn(user.Password)

	if err != nil {
		return entities.Session{}, err
	}

	user.Password = hashed
	user.Status = "NO_ORG"

	user, err = service.usersService.Create(user)

	if err != nil {
		return entities.Session{}, err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authImpl) OrganizationSignUp(user entities.User) (entities.Session, error) {
	hashed, err := service.passwordHashFn(user.Password)

	if err != nil {
		return entities.Session{}, err
	}

	user.Password = hashed
	user.Status = "ORG_MEMBER"

	user, err = service.usersService.Create(user)

	if err != nil {
		return entities.Session{}, err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authImpl) SignIn(email, password string) (entities.Session, error) {
	user, err := service.usersService.FindOneByEmail(email)

	if err != nil {
		return entities.Session{}, err
	}

	if err = service.passwordCompareFn(password, user.Password); err != nil {
		return entities.Session{}, err
	}

	session, err := service.sessionsService.Generate(user.ID)

	if err != nil {
		return entities.Session{}, err
	}

	return session, nil
}

func (service *authImpl) SignOut(sessionId string) error {
	return service.sessionsService.DeleteOneBySessionId(sessionId)
}

func (service *authImpl) AuthenticateSession(sessionId string) (entities.User, error) {
	session, err := service.sessionsService.FindOneBySessionId(sessionId)

	if err != nil {
		return entities.User{}, err
	}

	user, err := service.usersService.FindOneById(session.UserID)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
