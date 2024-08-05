package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type HousingRooms interface {
	CreateMany(data []entities.HousingRoom, housingId string) ([]entities.HousingRoom, error)
	FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error)
	FindOneById(id string) (entities.HousingRoom, error)
	UpdateOneById(id string, data entities.HousingRoom) error
	IncreaseAvailableVacanciesById(id string) error
	DecreaseAvailableVacanciesById(id string) error
	InactivateOneById(id string) error
}

type housingRoomsImpl struct {
	repository repositories.HousingRooms
}

func NewHousingRooms(repository repositories.HousingRooms) HousingRooms {
	return &housingRoomsImpl{
		repository: repository,
	}
}

func (service *housingRoomsImpl) CreateMany(data []entities.HousingRoom, housingId string) ([]entities.HousingRoom, error) {
	return service.repository.CreateMany(data, housingId)
}

func (service *housingRoomsImpl) FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error) {
	return service.repository.FindManyByHousingId(housingId, limit, offset)
}

func (service *housingRoomsImpl) FindOneById(id string) (entities.HousingRoom, error) {
	return service.repository.FindOneById(id)
}

func (service *housingRoomsImpl) UpdateOneById(id string, data entities.HousingRoom) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *housingRoomsImpl) IncreaseAvailableVacanciesById(id string) error {
	return service.repository.IncreaseAvailableVacanciesById(id)
}

func (service *housingRoomsImpl) DecreaseAvailableVacanciesById(id string) error {
	return service.repository.DecreaseAvailableVacanciesById(id)
}

func (service *housingRoomsImpl) InactivateOneById(id string) error {
	room, err := service.FindOneById(id)

	if err != nil {
		return err
	}

	room.Status = utils.InactiveStatus

	return service.repository.UpdateOneById(id, room)
}
