package organizations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, data entities.Organization) (entities.Organization, error)
}

type createImpl struct {
	organizationsRepository repositories.Organizations
	usersRepository         repositories.Users
}

func NewCreate(
	organizationsRepository repositories.Organizations,
	usersRepository repositories.Users,
) Create {
	return &createImpl{
		organizationsRepository: organizationsRepository,
		usersRepository:         usersRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, data entities.Organization) (entities.Organization, error) {
	if err := guards.AuthorizeCreateOrganization(actor); err != nil {
		return entities.Organization{}, err
	}

	organization, err := uc.organizationsRepository.Create(data)

	if err != nil {
		return entities.Organization{}, err
	}

	actor.OrganizationID = organization.ID
	actor.PlatformRole = utils.OrgAdminPlatformRole

	if err = uc.usersRepository.UpdateOneByID(actor.ID, actor); err != nil {
		return entities.Organization{}, err
	}

	return organization, nil
}
