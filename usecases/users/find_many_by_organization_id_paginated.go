package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.User, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	usersRepository         repositories.Users
	organizationsRepository repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	usersRepository repositories.Users,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		usersRepository:         usersRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.User, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.usersRepository.FindManyByOrganizationIDPaginated(organizationID, offset, limit)
}
