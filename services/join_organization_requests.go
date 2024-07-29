package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type JoinOrganizationRequests interface {
	Create(userID string, data entities.JoinOrganizationRequest) (entities.JoinOrganizationRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	Accept(id string) error
	Reject(id string) error
}

type joinOrganizationRequestsImpl struct {
	usersService Users
	repository   repositories.JoinOrganizationRequests
}

func NewJoinOrganizationRequests(usersService Users, repository repositories.JoinOrganizationRequests) JoinOrganizationRequests {
	return &joinOrganizationRequestsImpl{
		usersService: usersService,
		repository:   repository,
	}
}

func (service *joinOrganizationRequestsImpl) Create(userID string, data entities.JoinOrganizationRequest) (entities.JoinOrganizationRequest, error) {
	data.UserID = userID
	return service.repository.Create(data)
}

func (service *joinOrganizationRequestsImpl) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	return service.repository.FindManyByOrganizationId(organizationId, offset, limit)
}

func (service *joinOrganizationRequestsImpl) Accept(id string) error {
	request, err := service.repository.FindOneAndDeleteById(id)

	if err != nil {
		return err
	}

	data := entities.User{
		OrganizationID: request.OrganizationID,
	}

	if err = service.usersService.UpdateOneById(request.UserID, data); err != nil {
		return err
	}

	return nil
}

func (service *joinOrganizationRequestsImpl) Reject(id string) error {
	return service.repository.DeleteOneById(id)
}
