package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Users interface {
	Create(data entities.User) (entities.User, error)
	FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.User, error)
	FindOneByID(id string) (entities.User, error)
	FindOneCompleteByID(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	UpdateOneByID(id string, data entities.User) error
	InactivateOneByID(id string) error
	ExistsByEmail(email string) (bool, error)
	ExistsByID(id string) (bool, error)
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

func (service *usersImpl) FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.User, error) {
	return service.repository.FindManyByOrganizationID(organizationID, offset, limit)
}

func (service *usersImpl) FindOneByID(id string) (entities.User, error) {
	return service.repository.FindOneByID(id)
}

func (service *usersImpl) FindOneCompleteByID(id string) (entities.User, error) {
	return service.repository.FindOneCompleteByID(id)
}

func (service *usersImpl) FindOneByEmail(email string) (entities.User, error) {
	return service.repository.FindOneByEmail(email)
}

func (service *usersImpl) UpdateOneByID(id string, data entities.User) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *usersImpl) InactivateOneByID(id string) error {
	user, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	user.Status = utils.InactiveStatus

	return service.repository.UpdateOneByID(user.ID, user)
}

func (service *usersImpl) ExistsByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (service *usersImpl) ExistsByID(id string) (bool, error) {
	count, err := service.repository.CountByID(id)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
