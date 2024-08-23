package join_platform_invites

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type ConsumeByCode interface {
	Execute(code string) (entities.JoinPlatformInvite, error)
}

type consumeByCodeImpl struct {
	joinPlatformInvitesRepository repositories.JoinPlatformInvites
}

func NewConsumeByCode(joinPlatformInvitesRepository repositories.JoinPlatformInvites) ConsumeByCode {
	return &consumeByCodeImpl{joinPlatformInvitesRepository: joinPlatformInvitesRepository}
}

func (uc *consumeByCodeImpl) Execute(code string) (entities.JoinPlatformInvite, error) {
	return uc.joinPlatformInvitesRepository.FindOneAndDeleteByCode(code)
}
