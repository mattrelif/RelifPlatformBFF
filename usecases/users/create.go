package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(data entities.User) (entities.User, error)
}

type create struct {
	usersRepository repositories.Users
}

func NewCreate(usersRepository repositories.Users) Create {
	return &create{
		usersRepository: usersRepository,
	}
}

func (uc *create) Execute(data entities.User) (entities.User, error) {
	count, err := uc.usersRepository.CountByEmail(data.Email)

	if err != nil {
		return entities.User{}, err
	}

	if count > 0 {
		return entities.User{}, utils.ErrUserAlreadyExists
	}

	return uc.usersRepository.Create(data)
}
