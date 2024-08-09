package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type HousingRooms interface {
	CreateMany(data []entities.HousingRoom, housingID string) ([]entities.HousingRoom, error)
	FindManyByHousingID(housingID string, limit, offset int64) (int64, []entities.HousingRoom, error)
	FindOneByID(id string) (entities.HousingRoom, error)
	FindOneCompleteByID(id string) (entities.HousingRoom, error)
	UpdateOneByID(id string, data entities.HousingRoom) error
	IncreaseAvailableVacanciesByID(id string) error
	DecreaseAvailableVacanciesByID(id string) error
	InactivateOneByID(id string) error
}

type housingRoomsImpl struct {
	repository repositories.HousingRooms
}

func NewHousingRooms(repository repositories.HousingRooms) HousingRooms {
	return &housingRoomsImpl{
		repository: repository,
	}
}

func (service *housingRoomsImpl) CreateMany(data []entities.HousingRoom, housingID string) ([]entities.HousingRoom, error) {
	return service.repository.CreateMany(data, housingID)
}

func (service *housingRoomsImpl) FindManyByHousingID(housingID string, limit, offset int64) (int64, []entities.HousingRoom, error) {
	return service.repository.FindManyByHousingID(housingID, limit, offset)
}

func (service *housingRoomsImpl) FindOneByID(id string) (entities.HousingRoom, error) {
	return service.repository.FindOneByID(id)
}

func (service *housingRoomsImpl) FindOneCompleteByID(id string) (entities.HousingRoom, error) {
	return service.repository.FindOneCompleteByID(id)
}

func (service *housingRoomsImpl) UpdateOneByID(id string, data entities.HousingRoom) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *housingRoomsImpl) IncreaseAvailableVacanciesByID(id string) error {
	return service.repository.IncreaseAvailableVacanciesByID(id)
}

func (service *housingRoomsImpl) DecreaseAvailableVacanciesByID(id string) error {
	return service.repository.DecreaseAvailableVacanciesByID(id)
}

func (service *housingRoomsImpl) InactivateOneByID(id string) error {
	room, err := service.FindOneByID(id)

	if err != nil {
		return err
	}

	room.Status = utils.InactiveStatus

	return service.repository.UpdateOneByID(id, room)
}
