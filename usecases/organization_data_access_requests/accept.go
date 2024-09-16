package organization_data_access_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Accept interface {
	Execute(actor entities.User, requestID string) error
}

type acceptImpl struct {
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests
	organizationDataAccessGrantsRepository   repositories.OrganizationDataAccessGrants
	organizationsRepository                  repositories.Organizations
}

func NewAccept(
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests,
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants,
	organizationsRepository repositories.Organizations,
) Accept {
	return &acceptImpl{
		organizationDataAccessRequestsRepository: organizationDataAccessRequestsRepository,
		organizationDataAccessGrantsRepository:   organizationDataAccessGrantsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *acceptImpl) Execute(actor entities.User, requestID string) error {
	request, err := uc.organizationDataAccessRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	targetOrganization, err := uc.organizationsRepository.FindOneByID(request.TargetOrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, targetOrganization); err != nil {
		return err
	}

	requesterOrganization, err := uc.organizationsRepository.FindOneByID(request.RequesterOrganizationID)

	if err != nil {
		return err
	}

	grant := entities.OrganizationDataAccessGrant{
		TargetOrganizationID: targetOrganization.ID,
		OrganizationID:       requesterOrganization.ID,
		AuditorID:            actor.ID,
	}

	if err = uc.organizationDataAccessGrantsRepository.Create(grant); err != nil {
		return err
	}

	requesterOrganization.AccessGrantedIDs = append(requesterOrganization.AccessGrantedIDs, targetOrganization.ID)

	if err = uc.organizationsRepository.UpdateOneByID(requesterOrganization.ID, requesterOrganization); err != nil {
		return err
	}

	request.Status = utils.AcceptedStatus
	request.AcceptedAt = time.Now()
	request.AuditorID = actor.ID

	if err = uc.organizationDataAccessRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
