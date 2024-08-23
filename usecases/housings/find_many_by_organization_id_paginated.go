package housings

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.Housing, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.Housing, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.housingsRepository.FindManyByOrganizationIDPaginated(organizationID, search, offset, limit)
}
