package organization_data_access_requests

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
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests
}

func NewReject(organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests) Reject {
	return &rejectImpl{
		organizationDataAccessRequestsRepository: organizationDataAccessRequestsRepository,
	}
}

func (uc *rejectImpl) Execute(actor entities.User, requestID, rejectReason string) error {
	if err := guards.IsSuperUser(actor); err != nil {
		return err
	}

	request, err := uc.organizationDataAccessRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	request.Status = utils.RejectedStatus
	request.RejectReason = rejectReason
	request.RejectedAt = time.Now()

	if err = uc.organizationDataAccessRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
