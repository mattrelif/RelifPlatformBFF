package users

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyRelifMembersPaginated interface {
	Execute(actor entities.User, offset, limit int64) (int64, []entities.User, error)
}

type findManyRelifMembersPaginatedImpl struct {
	usersRepository repositories.Users
}

func NewFindManyRelifMembersPaginated(usersRepository repositories.Users) FindManyRelifMembersPaginated {
	return &findManyRelifMembersPaginatedImpl{
		usersRepository: usersRepository,
	}
}

func (uc *findManyRelifMembersPaginatedImpl) Execute(actor entities.User, offset, limit int64) (int64, []entities.User, error) {
	if err := guards.IsSuperUser(actor); err != nil {
		return 0, nil, err
	}

	return uc.usersRepository.FindManyRelifMembersPaginated(offset, limit)
}
