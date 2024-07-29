package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type Users interface {
	Create(data entities.User) (entities.User, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	FindOneAndUpdateById(id string, data entities.User) (entities.User, error)
	UpdateOneById(id string, data entities.User) error
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

func (service *usersImpl) Create(data entities.User) (entities.User, error) {
	data.Status = "ACTIVE"
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

func (service *usersImpl) FindOneAndUpdateById(id string, data entities.User) (entities.User, error) {
	return service.repository.FindOneAndUpdateById(id, data)
}

func (service *usersImpl) UpdateOneById(id string, data entities.User) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *usersImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
