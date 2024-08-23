package voluntary_people

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, id string) error
}

type deleteOneByID struct {
	voluntaryPeopleRepository repositories.VoluntaryPeople
	organizationsRepository   repositories.Organizations
}

func NewDeleteOneByID(
	voluntaryPeopleRepository repositories.VoluntaryPeople,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByID{
		voluntaryPeopleRepository: voluntaryPeopleRepository,
		organizationsRepository:   organizationsRepository,
	}
}

func (uc *deleteOneByID) Execute(actor entities.User, id string) error {
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

	return uc.voluntaryPeopleRepository.DeleteOneByID(voluntary.ID)
}
