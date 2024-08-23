package organizations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type ReactivateOneByID interface {
	Execute(actor entities.User, id string) error
}

type reactivateOneByIDImpl struct {
	organizationsRepository repositories.Organizations
}

func NewReactivateOneByID(organizationsRepository repositories.Organizations) ReactivateOneByID {
	return &reactivateOneByIDImpl{
		organizationsRepository: organizationsRepository,
	}
}

func (uc *reactivateOneByIDImpl) Execute(actor entities.User, id string) error {
	organization, err := uc.organizationsRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	if err = guards.IsSuperUser(actor); err != nil {
		return err
	}

	organization.Status = utils.ActiveStatus

	if err = uc.organizationsRepository.UpdateOneByID(id, organization); err != nil {
		return err
	}

	return nil
}
