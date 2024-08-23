package voluntary_people

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindOneByID interface {
	Execute(actor entities.User, id string) (entities.VoluntaryPerson, error)
}

type findOneByIDImpl struct {
	voluntaryPeopleRepository repositories.VoluntaryPeople
	organizationsRepository   repositories.Organizations
}

func NewFindOneByID(
	voluntaryPeopleRepository repositories.VoluntaryPeople,
	organizationsRepository repositories.Organizations,
) FindOneByID {
	return &findOneByIDImpl{
		voluntaryPeopleRepository: voluntaryPeopleRepository,
		organizationsRepository:   organizationsRepository,
	}
}

func (uc *findOneByIDImpl) Execute(actor entities.User, id string) (entities.VoluntaryPerson, error) {
	voluntary, err := uc.voluntaryPeopleRepository.FindOneByID(id)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(voluntary.OrganizationID)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	return voluntary, nil
}
