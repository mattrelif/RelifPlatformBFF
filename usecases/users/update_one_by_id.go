package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.User) error
}

type updateOneByIDImpl struct {
	usersRepository repositories.Users
}

func NewUpdateOneByID(usersRepository repositories.Users) UpdateOneByID {
	return &updateOneByIDImpl{
		usersRepository: usersRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.User) error {
	user, err := uc.usersRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	if err = guards.IsUser(actor, user); err != nil {
		return err
	}

	if err = uc.usersRepository.UpdateOneByID(id, data); err != nil {
		return err
	}

	return nil
}
