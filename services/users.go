package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Users interface {
	Create(data entities.User) (entities.User, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	UpdateOneById(id string, data entities.User) error
	InactivateOneById(id string) error
	AuthorizeExternalMutation(id string, user entities.User) error
	ExistsByEmail(email string) (bool, error)
	ExistsById(id string) (bool, error)
}

type usersImpl struct {
	repository repositories.Users
}

func NewUsers(repository repositories.Users) Users {
	return &usersImpl{
		repository: repository,
	}
}

func (service *usersImpl) Create(data entities.User) (entities.User, error) {
	exists, err := service.ExistsByEmail(data.Email)

	if err != nil {
		return entities.User{}, err
	}

	if exists {
		return entities.User{}, utils.ErrUserAlreadyExists
	}

	return service.repository.CreateUser(data)
}

func (service *usersImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *usersImpl) FindOneById(id string) (entities.User, error) {
	return service.repository.FindOneById(id)
}

func (service *usersImpl) FindOneByEmail(email string) (entities.User, error) {
	return service.repository.FindOneByEmail(email)
}

func (service *usersImpl) UpdateOneById(id string, data entities.User) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *usersImpl) InactivateOneById(id string) error {
	return service.repository.UpdateOneById(id, entities.User{Status: utils.InactiveStatus})
}

func (service *usersImpl) AuthorizeExternalMutation(id string, user entities.User) error {
	target, err := service.FindOneById(id)

	if err != nil {
		return err
	}

	if target.ID != user.ID && (target.OrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *usersImpl) ExistsByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (service *usersImpl) ExistsById(id string) (bool, error) {
	count, err := service.repository.CountById(id)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
