package join_platform_admin_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyPaginated interface {
	Execute(actor entities.User, offset, limit int64) (int64, []entities.JoinPlatformAdminInvite, error)
}

type findManyPaginatedImpl struct {
	repository repositories.JoinPlatformAdminInvites
}

func NewFindManyPaginated(repository repositories.JoinPlatformAdminInvites) FindManyPaginated {
	return &findManyPaginatedImpl{
		repository: repository,
	}
}

func (uc *findManyPaginatedImpl) Execute(actor entities.User, offset, limit int64) (int64, []entities.JoinPlatformAdminInvite, error) {
	if err := guards.IsSuperUser(actor); err != nil {
		return 0, nil, err
	}

	return uc.repository.FindManyPaginated(offset, limit)
}
