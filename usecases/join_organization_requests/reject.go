package join_organization_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Reject interface {
	Execute(actor entities.User, requestID, rejectReason string) error
}

type rejectImpl struct {
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests
	organizationsRepository            repositories.Organizations
}

func NewReject(
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests,
	organizationsRepository repositories.Organizations,
) Reject {
	return &rejectImpl{
		joinOrganizationRequestsRepository: joinOrganizationRequestsRepository,
		organizationsRepository:            organizationsRepository,
	}
}

func (uc *rejectImpl) Execute(actor entities.User, requestID, rejectReason string) error {
	request, err := uc.joinOrganizationRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(request.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	request.Status = utils.RejectedStatus
	request.RejectReason = rejectReason
	request.RejectedAt = time.Now()

	if err = uc.joinOrganizationRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
