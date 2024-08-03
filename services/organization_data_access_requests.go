package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type OrganizationDataAccessRequests interface {
	Create(requester entities.User, targetOrganizationId string) (entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindOneById(id string) (entities.OrganizationDataAccessRequest, error)
	Accept(id string, auditor entities.User) error
	Reject(id string, auditor entities.User, data entities.OrganizationDataAccessRequest) error
	AuthorizeCreate(user entities.User) error
	AuthorizeFindManyByOrganizationId(user entities.User, organizationId string) error
	AuthorizeExternalMutation(user entities.User, id string) (entities.OrganizationDataAccessRequest, error)
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

func (service *accessOrganizationDataRequestsImpl) Create(requester entities.User, targetOrganizationId string) (entities.OrganizationDataAccessRequest, error) {
	data := entities.OrganizationDataAccessRequest{
		RequesterID:             requester.ID,
		RequesterOrganizationID: requester.OrganizationID,
		TargetOrganizationID:    targetOrganizationId,
	}
	return service.repository.Create(data)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByRequesterOrganizationId(organizationId, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindManyByTargetOrganizationId(organizationId, limit, offset)
}

func (service *accessOrganizationDataRequestsImpl) FindOneById(id string) (entities.OrganizationDataAccessRequest, error) {
	return service.repository.FindOneById(id)
}

func (service *accessOrganizationDataRequestsImpl) Accept(id string, auditor entities.User) error {
	request, err := service.AuthorizeExternalMutation(auditor, id)

	if err != nil {
		return err
	}

	data := entities.OrganizationDataAccessRequest{
		AuditorID:  auditor.ID,
		Status:     utils.AcceptedStatus,
		AcceptedAt: time.Now(),
	}

	if err = service.repository.UpdateOneById(request.ID, data); err != nil {
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
	request, err := service.AuthorizeExternalMutation(auditor, id)

	if err != nil {
		return err
	}

	data.AuditorID = auditor.ID
	data.Status = utils.RejectedStatus
	data.RejectedAt = time.Now()

	if err = service.repository.UpdateOneById(request.ID, data); err != nil {
		return err
	}

	return nil
}

func (service *accessOrganizationDataRequestsImpl) AuthorizeCreate(user entities.User) error {
	organization, err := service.organizationsService.FindOneById(user.OrganizationID)

	if err != nil {
		return err
	}

	if organization.Type != utils.CoordinatorOrganizationType && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *accessOrganizationDataRequestsImpl) AuthorizeFindManyByOrganizationId(user entities.User, organizationId string) error {
	if user.OrganizationID != organizationId && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *accessOrganizationDataRequestsImpl) AuthorizeExternalMutation(user entities.User, id string) (entities.OrganizationDataAccessRequest, error) {
	request, err := service.repository.FindOneById(id)

	if err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	if request.RequesterOrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole {
		return entities.OrganizationDataAccessRequest{}, utils.ErrUnauthorizedAction
	}

	return request, nil
}
