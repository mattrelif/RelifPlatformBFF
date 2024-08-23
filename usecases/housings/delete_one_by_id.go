package housings

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, id string) error
}

type deleteOneByIDImpl struct {
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewDeleteOneByID(
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByIDImpl{
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *deleteOneByIDImpl) Execute(actor entities.User, id string) error {
	housing, err := uc.housingsRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	return uc.housingsRepository.DeleteOneByID(housing.ID)
}
