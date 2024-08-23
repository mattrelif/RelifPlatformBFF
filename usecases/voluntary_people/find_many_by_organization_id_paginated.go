package voluntary_people

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.VoluntaryPerson, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	voluntaryPeopleRepository repositories.VoluntaryPeople
	organizationsRepository   repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	voluntaryPeopleRepository repositories.VoluntaryPeople,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		voluntaryPeopleRepository: voluntaryPeopleRepository,
		organizationsRepository:   organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.VoluntaryPerson, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.voluntaryPeopleRepository.FindManyByOrganizationIDPaginated(organizationID, search, offset, limit)
}
