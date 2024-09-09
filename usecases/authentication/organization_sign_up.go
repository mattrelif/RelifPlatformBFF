package authentication

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	usersUseCase "relif/platform-bff/usecases/users"
	"relif/platform-bff/utils"
)

type OrganizationSignUp interface {
	Execute(data entities.User) (string, error)
}

type organizationSignUpImpl struct {
	sessionsRepository      repositories.Sessions
	organizationsRepository repositories.Organizations
	tokensService           services.Tokens
	createUserUseCase       usersUseCase.Create
	passwordHashFunction    utils.PasswordHashFn
}

func NewOrganizationSignUp(
	sessionsRepository repositories.Sessions,
	organizationsRepository repositories.Organizations,
	tokensService services.Tokens,
	createUserUseCase usersUseCase.Create,
	passwordHashFunction utils.PasswordHashFn,
) OrganizationSignUp {
	return &organizationSignUpImpl{
		sessionsRepository:      sessionsRepository,
		organizationsRepository: organizationsRepository,
		tokensService:           tokensService,
		createUserUseCase:       createUserUseCase,
		passwordHashFunction:    passwordHashFunction,
	}
}

func (uc *organizationSignUpImpl) Execute(data entities.User) (string, error) {
	organization, err := uc.organizationsRepository.FindOneByID(data.OrganizationID)

	if err != nil {
		return "", err
	}

	if organization.Status == utils.InactiveStatus {
		return "", utils.ErrMemberOfInactiveOrganization
	}

	hashed, err := uc.passwordHashFunction(data.Password)

	if err != nil {
		return "", err
	}

	data.Password = hashed
	data.PlatformRole = utils.OrgMemberPlatformRole
	data.OrganizationID = organization.ID

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
