package organizations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.Organization) error
}

type updateOneByIDImpl struct {
	organizationsRepository repositories.Organizations
}

func NewUpdateOneByID(organizationRepository repositories.Organizations) UpdateOneByID {
	return &updateOneByIDImpl{
		organizationsRepository: organizationRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.Organization) error {
	organization, err := uc.organizationsRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	if err = uc.organizationsRepository.UpdateOneByID(organization.ID, data); err != nil {
		return err
	}

	return nil
}
