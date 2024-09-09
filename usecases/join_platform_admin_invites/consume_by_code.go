package join_platform_admin_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type ConsumeByCode interface {
	Execute(code string) (entities.JoinPlatformAdminInvite, error)
}

type consumeByCode struct {
	repository repositories.JoinPlatformAdminInvites
}

func NewConsumeByCode(repository repositories.JoinPlatformAdminInvites) ConsumeByCode {
	return &consumeByCode{
		repository: repository,
	}
}

func (uc *consumeByCode) Execute(code string) (entities.JoinPlatformAdminInvite, error) {
	return uc.repository.FindOneAndDeleteByCode(code)
}
