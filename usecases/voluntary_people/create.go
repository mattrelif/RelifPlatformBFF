package voluntary_people

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, organizationID string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
}

type createImpl struct {
	voluntaryPeopleRepository repositories.VoluntaryPeople
	organizationsRepository   repositories.Organizations
}

func NewCreate(
	voluntaryPeopleRepository repositories.VoluntaryPeople,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		voluntaryPeopleRepository: voluntaryPeopleRepository,
		organizationsRepository:   organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, organizationID string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	count, err := uc.voluntaryPeopleRepository.CountByEmail(data.Email)

	if err != nil {
		return entities.VoluntaryPerson{}, err
	}

	if count > 0 {
		return entities.VoluntaryPerson{}, utils.ErrVoluntaryPersonAlreadyExists
	}

	data.OrganizationID = organization.ID

	return uc.voluntaryPeopleRepository.Create(data)
}
