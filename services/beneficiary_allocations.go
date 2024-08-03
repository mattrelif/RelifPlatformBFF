package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
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
	housingsService      Housings
}

func NewBeneficiaryAllocations(
	repository repositories.BeneficiaryAllocations,
	beneficiariesService Beneficiaries,
	housingRoomsService HousingRooms,
	housingsService Housings,
) BeneficiaryAllocations {
	return &beneficiaryAllocationsImpl{
		repository:           repository,
		beneficiariesService: beneficiariesService,
		housingRoomsService:  housingRoomsService,
		housingsService:      housingsService,
	}
}

func (service *beneficiaryAllocationsImpl) Allocate(user entities.User, beneficiaryId string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	data.AuditorID = user.ID
	data.Type = utils.EntranceType
	data.BeneficiaryID = beneficiaryId

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	housing, err := service.housingsService.FindOneByID(allocation.HousingID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.beneficiariesService.UpdateOneById(allocation.BeneficiaryID, entities.Beneficiary{CurrentOrganizationID: housing.OrganizationID, CurrentRoomID: allocation.RoomID, CurrentHousingID: allocation.HousingID}); err != nil {
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
	data.Type = utils.ReallocationType
	data.OldRoomID = beneficiary.CurrentRoomID
	data.OldHousingID = beneficiary.CurrentHousingID
	data.BeneficiaryID = beneficiaryId

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	housing, err := service.housingsService.FindOneByID(allocation.HousingID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = service.beneficiariesService.UpdateOneById(allocation.BeneficiaryID, entities.Beneficiary{CurrentOrganizationID: housing.OrganizationID, CurrentRoomID: allocation.RoomID, CurrentHousingID: allocation.HousingID}); err != nil {
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
