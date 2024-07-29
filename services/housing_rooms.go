package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type HousingRooms interface {
	CreateMany(data []entities.HousingRoom) ([]entities.HousingRoom, error)
	FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error)
	FindOneById(id string) (entities.HousingRoom, error)
	FindOneAndUpdateById(id string, data entities.HousingRoom) (entities.HousingRoom, error)
	IncreaseAvailableVacanciesById(id string) error
	DecreaseAvailableVacanciesById(id string) error
	DeleteOneById(id string) error
}

type housingRoomsImpl struct {
	repository repositories.HousingRooms
}

func NewHousingRooms(repository repositories.HousingRooms) HousingRooms {
	return &housingRoomsImpl{
		repository: repository,
	}
}

func (service *housingRoomsImpl) CreateMany(data []entities.HousingRoom) ([]entities.HousingRoom, error) {
	return service.repository.CreateMany(data)
}

func (service *housingRoomsImpl) FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error) {
	return service.repository.FindManyByHousingId(housingId, limit, offset)
}

func (service *housingRoomsImpl) FindOneById(id string) (entities.HousingRoom, error) {
	return service.repository.FindOneById(id)
}

func (service *housingRoomsImpl) FindOneAndUpdateById(id string, data entities.HousingRoom) (entities.HousingRoom, error) {
	return service.repository.FindOneAndUpdateById(id, data)
}

func (service *housingRoomsImpl) IncreaseAvailableVacanciesById(id string) error {
	return service.repository.IncreaseAvailableVacanciesById(id)
}

func (service *housingRoomsImpl) DecreaseAvailableVacanciesById(id string) error {
	return service.repository.DecreaseAvailableVacanciesById(id)
}

func (service *housingRoomsImpl) DeleteOneById(id string) error {
	return service.repository.DeleteOneById(id)
}
