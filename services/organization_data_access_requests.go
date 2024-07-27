package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type OrganizationDataAccessRequests interface {
	Create(requester entities.User, request entities.OrganizationDataAccessRequest) (string, error)
	FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	Accept(id string, userId string) error
	Reject(id string, userId string, request entities.OrganizationDataAccessRequest) error
}

type accessOrganizationDataRequestsImpl struct {
	repository    repositories.OrganizationDataAccessRequests
	grantsService OrganizationDataAccessGrants
}

func NewOrganizationDataAccessRequests(
	repository repositories.OrganizationDataAccessRequests,
	grantsService OrganizationDataAccessGrants,
) OrganizationDataAccessRequests {
	return &accessOrganizationDataRequestsImpl{
		repository:    repository,
		grantsService: grantsService,
	}
}

func (service *accessOrganizationDataRequestsImpl) Create(requester entities.User, request entities.OrganizationDataAccessRequest) (string, error) {
	request.RequesterID = requester.ID
	request.RequesterOrganizationID = requester.OrganizationID
	request.CreatedAt = time.Now()

	return service.repository.Create(request)
}

func (service *accessOrganizationDataRequestsImpl) FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindMany(limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByRequesterOrganizationId(organizationId, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) Accept(id string, userId string) error {
	request := entities.OrganizationDataAccessRequest{
		AuditorID:  userId,
		Status:     "ACCEPTED",
		AcceptedAt: time.Now(),
	}

	updated, err := service.repository.FindOneAndUpdateById(id, request)

	if err != nil {
		return err
	}

	grant := entities.OrganizationDataAccessGrant{
		OrganizationID:       updated.RequesterOrganizationID,
		AuditorID:            updated.AuditorID,
		TargetOrganizationID: updated.TargetOrganizationID,
		CreatedAt:            time.Now(),
	}

	if err = service.grantsService.Create(grant); err != nil {
		return err
	}

	return nil
}

func (service *accessOrganizationDataRequestsImpl) Reject(id string, userId string, request entities.OrganizationDataAccessRequest) error {
	request.RejectedAt = time.Now()
	request.AuditorID = userId
	request.Status = "REJECTED"

	if err := service.repository.UpdateOneById(id, request); err != nil {
		return err
	}

	return nil
}
