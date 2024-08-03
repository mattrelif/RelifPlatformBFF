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
	AuthorizeFindMany(user entities.User) error
	AuthorizeFindManyByOrganizationId(user entities.User, organizationId string) error
	AuthorizeExternalMutation(user entities.User) error
	AuthorizeCreate(user entities.User) error
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

	if err = service.repository.UpdateOneById(id, entities.UpdateOrganizationTypeRequest{AuditorID: userId, AcceptedAt: time.Now(), Status: utils.AcceptedStatus}); err != nil {
		return err
	}

	if err = service.organizationsService.UpdateOneById(request.OrganizationID, entities.Organization{Type: utils.CoordinatorOrganizationType}); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) Reject(id, userId string, data entities.UpdateOrganizationTypeRequest) error {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	if err = service.repository.UpdateOneById(request.ID, entities.UpdateOrganizationTypeRequest{AuditorID: userId, Status: utils.RejectedStatus, RejectedAt: time.Now(), RejectReason: data.RejectReason}); err != nil {
		return err
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) AuthorizeFindMany(user entities.User) error {
	if user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) AuthorizeFindManyByOrganizationId(user entities.User, organizationId string) error {
	if user.OrganizationID != organizationId && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) AuthorizeExternalMutation(user entities.User) error {
	if user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *updateOrganizationTypeRequestsImpl) AuthorizeCreate(user entities.User) error {
	if user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	exists, err := service.ExistsPendingByOrganization(user.OrganizationID)

	if err != nil {
		return err
	}

	if exists {
		return utils.ErrUnauthorizedAction
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
