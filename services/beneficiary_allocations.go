package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
)

type BeneficiaryAllocations interface {
	Allocate(user entities.User, beneficiaryId string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
	Reallocate(user entities.User, beneficiaryId string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
	FindManyByBeneficiaryId(beneficiaryId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByHousingId(housingId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByRoomId(roomId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
}

type beneficiaryAllocationsImpl struct {
	repository           repositories.BeneficiaryAllocations
	beneficiariesService Beneficiaries
	housingRoomsService  HousingRooms
}

func NewBeneficiaryAllocations(
	repository repositories.BeneficiaryAllocations,
	beneficiariesService Beneficiaries,
	housingRoomsService HousingRooms,
) BeneficiaryAllocations {
	return &beneficiaryAllocationsImpl{
		repository:           repository,
		beneficiariesService: beneficiariesService,
		housingRoomsService:  housingRoomsService,
	}
}

func (service *beneficiaryAllocationsImpl) Allocate(user entities.User, beneficiaryId string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	data.AuditorID = user.ID
	data.Type = "ENTRANCE"
	data.BeneficiaryID = beneficiaryId

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.beneficiariesService.UpdateOneById(allocation.BeneficiaryID, entities.Beneficiary{CurrentRoomID: allocation.RoomID, CurrentHousingID: allocation.HousingID}); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.housingRoomsService.DecreaseAvailableVacanciesById(allocation.RoomID); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	return allocation, nil
}

func (service *beneficiaryAllocationsImpl) Reallocate(user entities.User, beneficiaryId string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	beneficiary, err := service.beneficiariesService.FindOneById(beneficiaryId)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	data.AuditorID = user.ID
	data.Type = "REALLOCATION"
	data.OldRoomID = beneficiary.CurrentRoomID
	data.OldHousingID = beneficiary.CurrentHousingID
	data.BeneficiaryID = beneficiaryId

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.beneficiariesService.UpdateOneById(allocation.BeneficiaryID, entities.Beneficiary{CurrentRoomID: allocation.RoomID, CurrentHousingID: allocation.HousingID}); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.housingRoomsService.IncreaseAvailableVacanciesById(allocation.OldRoomID); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.housingRoomsService.DecreaseAvailableVacanciesById(allocation.RoomID); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	return allocation, err
}

func (service *beneficiaryAllocationsImpl) FindManyByBeneficiaryId(beneficiaryId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByBeneficiaryId(beneficiaryId, offset, limit)
}

func (service *beneficiaryAllocationsImpl) FindManyByHousingId(housingId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByHousingId(housingId, offset, limit)
}

func (service *beneficiaryAllocationsImpl) FindManyByRoomId(roomId string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByRoomId(roomId, offset, limit)
}
