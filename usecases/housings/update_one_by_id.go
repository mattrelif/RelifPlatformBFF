package housings

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.Housing) error
}

type updateOneByIDImpl struct {
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewUpdateOneByID(
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) UpdateOneByID {
	return &updateOneByIDImpl{
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.Housing) error {
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

	return uc.housingsRepository.UpdateOneByID(housing.ID, data)
}
