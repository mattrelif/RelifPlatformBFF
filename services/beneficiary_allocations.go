package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type BeneficiaryAllocations interface {
	Allocate(user entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
	Reallocate(user entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
	FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByHousingID(housingID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByRoomID(roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
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

func (service *beneficiaryAllocationsImpl) Allocate(user entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	var room entities.HousingRoom

	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	housing, err := service.housingsService.FindOneByID(data.HousingID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if data.RoomID != "" {
		room, err = service.housingRoomsService.FindOneByID(data.RoomID)

		if err != nil {
			return entities.BeneficiaryAllocation{}, err
		}
	}

	data.AuditorID = user.ID
	data.Type = utils.EntranceType
	data.BeneficiaryID = beneficiaryID

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	beneficiary.CurrentHousingID = housing.ID
	beneficiary.CurrentOrganizationID = housing.OrganizationID
	beneficiary.CurrentRoomID = room.ID

	if err = service.beneficiariesService.UpdateOneByID(beneficiary.ID, beneficiary); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	return allocation, nil
}

func (service *beneficiaryAllocationsImpl) Reallocate(user entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	var room entities.HousingRoom

	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	housing, err := service.housingsService.FindOneByID(data.HousingID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if data.RoomID != "" {
		room, err = service.housingRoomsService.FindOneByID(data.RoomID)

		if err != nil {
			return entities.BeneficiaryAllocation{}, err
		}
	}

	data.AuditorID = user.ID
	data.Type = utils.ReallocationType
	data.OldRoomID = beneficiary.CurrentRoomID
	data.OldHousingID = beneficiary.CurrentHousingID
	data.BeneficiaryID = beneficiary.ID

	allocation, err := service.repository.Create(data)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	beneficiary.CurrentHousingID = housing.ID
	beneficiary.CurrentOrganizationID = housing.OrganizationID
	beneficiary.CurrentRoomID = room.ID

	if err = service.beneficiariesService.UpdateOneByID(beneficiary.ID, beneficiary); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	return allocation, err
}

func (service *beneficiaryAllocationsImpl) FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByBeneficiaryID(beneficiaryID, offset, limit)
}

func (service *beneficiaryAllocationsImpl) FindManyByHousingID(housingID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByHousingID(housingID, offset, limit)
}

func (service *beneficiaryAllocationsImpl) FindManyByRoomID(roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	return service.repository.FindManyByRoomID(roomID, offset, limit)
}
