package join_organization_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByUserIDPaginated interface {
	Execute(actor entities.User, userID string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
}

type findManyByUserIDPaginatedImpl struct {
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites
	usersRepository                   repositories.Users
}

func NewFindManyByUserIDPaginated(
	joinOrganizationInvitesRepository repositories.JoinOrganizationInvites,
	usersRepository repositories.Users,
) FindManyByUserIDPaginated {
	return &findManyByUserIDPaginatedImpl{
		joinOrganizationInvitesRepository: joinOrganizationInvitesRepository,
		usersRepository:                   usersRepository,
	}
}

func (uc *findManyByUserIDPaginatedImpl) Execute(actor entities.User, userID string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	user, err := uc.usersRepository.FindOneByID(userID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsUser(actor, user); err != nil {
		return 0, nil, err
	}

	return uc.joinOrganizationInvitesRepository.FindManyByUserIDPaginated(userID, offset, limit)
}
