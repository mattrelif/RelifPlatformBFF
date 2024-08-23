package update_organization_type_requets

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
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests
}

func NewReject(updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests) Reject {
	return &rejectImpl{
		updateOrganizationTypeRequestsRepository: updateOrganizationTypeRequestsRepository,
	}
}

func (uc *rejectImpl) Execute(actor entities.User, requestID, rejectReason string) error {
	if err := guards.IsSuperUser(actor); err != nil {
		return err
	}

	request, err := uc.updateOrganizationTypeRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	request.Status = utils.RejectedStatus
	request.RejectReason = rejectReason
	request.RejectedAt = time.Now()

	if err = uc.updateOrganizationTypeRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
