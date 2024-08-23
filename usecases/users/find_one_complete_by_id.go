package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type FindOneCompleteByID interface {
	Execute(userID string) (entities.User, error)
}

type findOneCompleteByIDImpl struct {
	usersRepository repositories.Users
}

func NewFindOneCompleteByID(
	usersRepository repositories.Users,
) FindOneCompleteByID {
	return &findOneCompleteByIDImpl{
		usersRepository: usersRepository,
	}
}

func (uc *findOneCompleteByIDImpl) Execute(userID string) (entities.User, error) {
	user, err := uc.usersRepository.FindOneAndLookupByID(userID)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
