package beneficiary_allocations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByHousingRoomIDPaginated interface {
	Execute(actor entities.User, roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
}

type findManyByHousingRoomIDPaginatedImpl struct {
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations
	housingsRepository               repositories.Housings
	housingRoomsRepository           repositories.HousingRooms
	organizationsRepository          repositories.Organizations
}

func NewFindManyByHousingRoomIDPaginated(
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations,
	housingsRepository repositories.Housings,
	housingRoomsRepository repositories.HousingRooms,
	organizationsRepository repositories.Organizations,
) FindManyByHousingRoomIDPaginated {
	return &findManyByHousingRoomIDPaginatedImpl{
		beneficiaryAllocationsRepository: beneficiaryAllocationsRepository,
		housingsRepository:               housingsRepository,
		housingRoomsRepository:           housingRoomsRepository,
		organizationsRepository:          organizationsRepository,
	}
}

func (uc *findManyByHousingRoomIDPaginatedImpl) Execute(actor entities.User, roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	room, err := uc.housingRoomsRepository.FindOneByID(roomID)

	if err != nil {
		return 0, nil, err
	}

	housing, err := uc.housingsRepository.FindOneByID(room.HousingID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.beneficiaryAllocationsRepository.FindManyByRoomIDPaginated(room.ID, offset, limit)
}
