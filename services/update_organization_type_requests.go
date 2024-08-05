package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type UpdateOrganizationTypeRequests interface {
	Create(user entities.User) (entities.UpdateOrganizationTypeRequest, error)
	FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	Accept(id, userId string) error
	Reject(id, userId string, data entities.UpdateOrganizationTypeRequest) error
	ExistsPendingByOrganization(organizationId string) (bool, error)
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

func (service *updateOrganizationTypeRequestsImpl) Create(user entities.User) (entities.UpdateOrganizationTypeRequest, error) {
	data := entities.UpdateOrganizationTypeRequest{
		CreatorID:      user.ID,
		OrganizationID: user.OrganizationID,
	}
	return service.repository.Create(data)
}

func (service *updateOrganizationTypeRequestsImpl) FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	return service.repository.FindMany(offset, limit)
}

func (service *updateOrganizationTypeRequestsImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *updateOrganizationTypeRequestsImpl) Accept(id, userId string) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	organization, err := service.organizationsService.FindOneById(request.OrganizationID)

	if err != nil {
		return err
	}

	request.Status = utils.AcceptedStatus
	request.AuditorID = userId
	request.AcceptedAt = time.Now()

	if err = service.repository.UpdateOneById(request.ID, request); err != nil {
		return err
	}

	organization.Type = utils.CoordinatorOrganizationType

	if err = service.organizationsService.UpdateOneById(organization.ID, organization); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) Reject(id, userId string, data entities.UpdateOrganizationTypeRequest) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	request.AuditorID = userId
	request.Status = utils.RejectedStatus
	request.RejectedAt = time.Now()
	request.RejectReason = data.RejectReason

	if err = service.repository.UpdateOneById(request.ID, request); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) ExistsPendingByOrganization(organizationId string) (bool, error) {
	count, err := service.repository.CountPendingByOrganizationId(organizationId)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
