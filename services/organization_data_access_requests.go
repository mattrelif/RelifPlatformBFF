package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type OrganizationDataAccessRequests interface {
	Create(requester entities.User, data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error)
	FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	Accept(id string, userId string) error
	Reject(id string, userId string, data entities.OrganizationDataAccessRequest) error
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

func (service *accessOrganizationDataRequestsImpl) Create(requester entities.User, data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error) {
	data.RequesterID = requester.ID
	data.RequesterOrganizationID = requester.OrganizationID
	return service.repository.Create(data)
}

func (service *accessOrganizationDataRequestsImpl) FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindMany(limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByRequesterOrganizationId(organizationId, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) Accept(id string, userId string) error {
	data := entities.OrganizationDataAccessRequest{
		AuditorID:  userId,
		Status:     "ACCEPTED",
		AcceptedAt: time.Now(),
	}

	updated, err := service.repository.FindOneAndUpdateById(id, data)

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

func (service *accessOrganizationDataRequestsImpl) Reject(id string, userId string, data entities.OrganizationDataAccessRequest) error {
	data.RejectedAt = time.Now()
	data.AuditorID = userId
	data.Status = "REJECTED"

	if err := service.repository.UpdateOneById(id, data); err != nil {
		return err
	}

	return nil
}
