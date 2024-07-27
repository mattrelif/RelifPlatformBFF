package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type UpdateOrganizationTypeRequests interface {
	Create(userId string, request entities.UpdateOrganizationTypeRequest) (string, error)
	FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	Accept(id, userId string) error
	Reject(id, userId string, request entities.UpdateOrganizationTypeRequest) error
}

type updateOrganizationTypeRequestsImpl struct {
	organizationsService Organizations
	repository           repositories.UpdateOrganizationTypeRequests
}

func NewUpdateOrganizationTypeRequests(
	organizationsService Organizations,
	repository repositories.UpdateOrganizationTypeRequests,
) UpdateOrganizationTypeRequests {
	return &updateOrganizationTypeRequestsImpl{
		organizationsService: organizationsService,
		repository:           repository,
	}
}

func (service *updateOrganizationTypeRequestsImpl) Create(userId string, request entities.UpdateOrganizationTypeRequest) (string, error) {
	request.CreatorID = userId
	request.Status = "PENDING"
	request.CreatedAt = time.Now()

	return service.repository.Create(request)
}

func (service *updateOrganizationTypeRequestsImpl) FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *updateOrganizationTypeRequestsImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *updateOrganizationTypeRequestsImpl) Accept(id, userId string) error {
	request, err := service.repository.FindOneAndUpdateById(id, entities.UpdateOrganizationTypeRequest{AuditorID: userId, Status: "ACCEPTED"})

	if err != nil {
		return err
	}

	if err = service.organizationsService.UpdateOneById(request.OrganizationID, entities.Organization{Type: "MANAGER"}); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) Reject(id, userId string, request entities.UpdateOrganizationTypeRequest) error {
	request, err := service.repository.FindOneAndUpdateById(id, entities.UpdateOrganizationTypeRequest{AuditorID: userId, Status: "REJECTED", RejectedAt: time.Now(), RejectReason: request.RejectReason})

	if err != nil {
		return err
	}

	return nil
}
