package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type Beneficiaries interface {
	Create(data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingId(housingId string) ([]entities.Beneficiary, error)
	FindManyByRoomId(roomId string) ([]entities.Beneficiary, error)
	FindOneById(id string) (entities.Beneficiary, error)
	FindOneAndUpdateById(id string, data entities.Beneficiary) (entities.Beneficiary, error)
	UpdateOneById(id string, data entities.Beneficiary) error
	DeleteOneById(id string) error
}

type beneficiariesImpl struct {
	repository repositories.Beneficiaries
}

func NewBeneficiaries(repository repositories.Beneficiaries) Beneficiaries {
	return &beneficiariesImpl{
		repository: repository,
	}
}

func (service *beneficiariesImpl) Create(data entities.Beneficiary) (entities.Beneficiary, error) {
	return service.repository.Create(data)
}

func (service *beneficiariesImpl) FindManyByHousingId(housingId string) ([]entities.Beneficiary, error) {
	return service.repository.FindManyByHousingId(housingId)
}

func (service *beneficiariesImpl) FindManyByRoomId(roomId string) ([]entities.Beneficiary, error) {
	return service.repository.FindManyByRoomId(roomId)
}

func (service *beneficiariesImpl) FindOneById(id string) (entities.Beneficiary, error) {
	return service.repository.FindOneById(id)
}

func (service *beneficiariesImpl) FindOneAndUpdateById(id string, data entities.Beneficiary) (entities.Beneficiary, error) {
	return service.repository.FindOneAndUpdateById(id, data)
}

func (service *beneficiariesImpl) UpdateOneById(id string, data entities.Beneficiary) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *beneficiariesImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
