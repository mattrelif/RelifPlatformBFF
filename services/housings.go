package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Housings interface {
	Create(user entities.User, housing entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string, user entities.User) (entities.Housing, error)
	UpdateOneByID(id string, housing entities.Housing) error
	InactivateOneByID(id string) error
	AuthorizeCreate(user entities.User) error
	AuthorizeFindManyByOrganizationID(user entities.User, organizationId string) error
	AuthorizeExternalMutation(user entities.User, id string) error
}

type housingsImpl struct {
	repository    repositories.Housings
	grantsService OrganizationDataAccessGrants
}

func NewHousings(repository repositories.Housings, grantsService OrganizationDataAccessGrants) Housings {
	return &housingsImpl{
		repository:    repository,
		grantsService: grantsService,
	}
}

func (service *housingsImpl) Create(user entities.User, housing entities.Housing) (entities.Housing, error) {
	housing.OrganizationID = user.OrganizationID
	return service.repository.Create(housing)
}

func (service *housingsImpl) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	return service.repository.FindManyByOrganizationID(organizationId, limit, offset)
}

func (service *housingsImpl) FindOneByID(id string, user entities.User) (entities.Housing, error) {
	return service.authorizeFindOneByID(user, id)
}

func (service *housingsImpl) UpdateOneByID(id string, housing entities.Housing) error {
	return service.repository.UpdateOneById(id, housing)
}

func (service *housingsImpl) InactivateOneByID(id string) error {
	data := entities.Housing{
		Status: utils.InactiveStatus,
	}
	return service.repository.UpdateOneById(id, data)
}

func (service *housingsImpl) AuthorizeCreate(user entities.User) error {
	if user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *housingsImpl) AuthorizeExternalMutation(user entities.User, id string) error {
	housing, err := service.repository.FindOneByID(id)

	if err != nil {
		return err
	}

	if housing.OrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *housingsImpl) authorizeFindOneByID(user entities.User, id string) (entities.Housing, error) {
	housing, err := service.repository.FindOneByID(id)

	if err != nil {
		return entities.Housing{}, err
	}

	accessGranted, err := service.grantsService.ExistsByOrganizationIdAndTargetOrganizationId(user.OrganizationID, housing.OrganizationID)

	if err != nil {
		return entities.Housing{}, err
	}

	if user.OrganizationID != housing.OrganizationID && !accessGranted {
		return entities.Housing{}, utils.ErrUnauthorizedAction
	}

	return housing, nil
}

func (service *housingsImpl) AuthorizeFindManyByOrganizationID(user entities.User, organizationId string) error {
	accessGranted, err := service.grantsService.ExistsByOrganizationIdAndTargetOrganizationId(user.OrganizationID, organizationId)

	if err != nil {
		return err
	}

	if user.OrganizationID != organizationId && !accessGranted {
		return utils.ErrUnauthorizedAction
	}

	return nil
}
