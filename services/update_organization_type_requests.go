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
	FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	Accept(id, userID string) error
	Reject(id, userID string, data entities.UpdateOrganizationTypeRequest) error
	ExistsPendingByOrganization(organizationID string) (bool, error)
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

func (service *updateOrganizationTypeRequestsImpl) FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	return service.repository.FindManyByOrganizationID(organizationID, offset, limit)
}

func (service *updateOrganizationTypeRequestsImpl) Accept(id, userID string) error {
	request, err := service.repository.FindOneByID(id)

	if err != nil {
		return err
	}

	organization, err := service.organizationsService.FindOneByID(request.OrganizationID)

	if err != nil {
		return err
	}

	request.Status = utils.AcceptedStatus
	request.AuditorID = userID
	request.AcceptedAt = time.Now()

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	organization.Type = utils.CoordinatorOrganizationType

	if err = service.organizationsService.UpdateOneByID(organization.ID, organization); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) Reject(id, userID string, data entities.UpdateOrganizationTypeRequest) error {
	request, err := service.repository.FindOneByID(id)

	if err != nil {
		return err
	}

	request.AuditorID = userID
	request.Status = utils.RejectedStatus
	request.RejectedAt = time.Now()
	request.RejectReason = data.RejectReason

	if err = service.repository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) ExistsPendingByOrganization(organizationID string) (bool, error) {
	count, err := service.repository.CountPendingByOrganizationID(organizationID)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
