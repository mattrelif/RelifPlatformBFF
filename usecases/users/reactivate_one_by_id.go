package users

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
	usersRepository repositories.Users
}

func NewReactivateOneByID(usersRepository repositories.Users) ReactivateOneByID {
	return &reactivateOneByIDImpl{
		usersRepository: usersRepository,
	}
}

func (uc *reactivateOneByIDImpl) Execute(actor entities.User, id string) error {
	user, err := uc.usersRepository.FindOneAndLookupByID(id)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, user.Organization); err != nil {
		return err
	}

	user.Status = utils.ActiveStatus

	if err = uc.usersRepository.UpdateOneByID(id, user); err != nil {
		return err
	}

	return nil
}
