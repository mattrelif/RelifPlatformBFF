package voluntary_people

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.VoluntaryPerson) error
}

type updateOneByIDImpl struct {
	voluntaryPeopleRepository repositories.VoluntaryPeople
	organizationsRepository   repositories.Organizations
}

func NewUpdateOneByID(
	voluntaryPeopleRepository repositories.VoluntaryPeople,
	organizationsRepository repositories.Organizations,
) UpdateOneByID {
	return &updateOneByIDImpl{
		voluntaryPeopleRepository: voluntaryPeopleRepository,
		organizationsRepository:   organizationsRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.VoluntaryPerson) error {
	voluntary, err := uc.voluntaryPeopleRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(voluntary.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	return uc.voluntaryPeopleRepository.UpdateOneByID(voluntary.ID, data)
}
