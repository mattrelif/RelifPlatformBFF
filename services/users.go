package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"time"
)

type Users interface {
	Create(user entities.User) (string, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	FindOneAndUpdateById(id string, user entities.User) (entities.User, error)
	UpdateOneById(id string, user entities.User) error
	DeleteOneById(id string) error
}

type usersImpl struct {
	repository repositories.Users
}

func NewUsers(repository repositories.Users) Users {
	return &usersImpl{
		repository: repository,
	}
}

func (service *usersImpl) Create(user entities.User) (string, error) {
	user.CreatedAt = time.Now()
	return service.repository.CreateUser(user)
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

func (service *usersImpl) FindOneAndUpdateById(id string, user entities.User) (entities.User, error) {
	user.UpdatedAt = time.Now()
	return service.repository.FindOneAndUpdateById(id, user)
}

func (service *usersImpl) UpdateOneById(id string, user entities.User) error {
	user.UpdatedAt = time.Now()
	return service.repository.UpdateOneById(id, user)
}

func (service *usersImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
