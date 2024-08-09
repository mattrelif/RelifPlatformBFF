package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type OrganizationDataAccessRequests interface {
	Create(requester entities.User, targetOrganizationID string) (entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByTargetOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindOneByID(id string) (entities.OrganizationDataAccessRequest, error)
	Accept(id string, auditor entities.User) error
	Reject(id string, auditor entities.User, data entities.OrganizationDataAccessRequest) error
}

type accessOrganizationDataRequestsImpl struct {
	repository           repositories.OrganizationDataAccessRequests
	organizationsService Organizations
	grantsService        OrganizationDataAccessGrants
}

func NewOrganizationDataAccessRequests(
	repository repositories.OrganizationDataAccessRequests,
	organizationsService Organizations,
	grantsService OrganizationDataAccessGrants,
) OrganizationDataAccessRequests {
	return &accessOrganizationDataRequestsImpl{
		repository:           repository,
		organizationsService: organizationsService,
		grantsService:        grantsService,
	}
}

func (service *accessOrganizationDataRequestsImpl) Create(requester entities.User, targetOrganizationID string) (entities.OrganizationDataAccessRequest, error) {
	if requester.OrganizationID == targetOrganizationID {
		return entities.OrganizationDataAccessRequest{}, utils.ErrUnauthorizedAction
	}

	data := entities.OrganizationDataAccessRequest{
		RequesterID:             requester.ID,
		RequesterOrganizationID: requester.OrganizationID,
		TargetOrganizationID:    targetOrganizationID,
	}

	return service.repository.Create(data)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByRequesterOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByRequesterOrganizationID(organizationID, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByTargetOrganizationID(organizationID string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByTargetOrganizationID(organizationID, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindOneByID(id string) (entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindOneByID(id)
}

func (service *accessOrganizationDataRequestsImpl) Accept(id string, auditor entities.User) error {
	request, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	request.AcceptedAt = time.Now()
	request.Status = utils.AcceptedStatus
	request.AuditorID = auditor.ID

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	grant := entities.OrganizationDataAccessGrant{
		OrganizationID:       request.RequesterOrganizationID,
		AuditorID:            request.AuditorID,
		TargetOrganizationID: request.TargetOrganizationID,
		CreatedAt:            time.Now(),
	}

	if err = service.grantsService.Create(grant); err != nil {
		return err
	}

	return nil
}

func (service *accessOrganizationDataRequestsImpl) Reject(id string, auditor entities.User, data entities.OrganizationDataAccessRequest) error {
	request, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	request.AuditorID = auditor.ID
	request.Status = utils.RejectedStatus
	request.RejectedAt = time.Now()
	request.RejectReason = data.RejectReason

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
