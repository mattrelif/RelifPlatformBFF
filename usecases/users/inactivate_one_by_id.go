package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type InactivateOneByID interface {
	Execute(actor entities.User, id string) error
}

type inactivateOneByIDImpl struct {
	usersRepository repositories.Users
}

func NewInactivateOneByID(usersRepository repositories.Users) InactivateOneByID {
	return &inactivateOneByIDImpl{
		usersRepository: usersRepository,
	}
}

func (uc *inactivateOneByIDImpl) Execute(actor entities.User, id string) error {
	user, err := uc.usersRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	if err = guards.IsSuperUser(actor); err != nil {
		return err
	}

	user.Status = utils.InactiveStatus

	if err = uc.usersRepository.UpdateOneByID(id, user); err != nil {
		return err
	}

	return nil
}
