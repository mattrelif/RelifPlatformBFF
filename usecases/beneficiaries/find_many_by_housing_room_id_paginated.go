package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByHousingRoomIDPaginated interface {
	Execute(actor entities.User, roomID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
}

type findManyByHousingRoomIDPaginatedImpl struct {
	beneficiariesRepository repositories.Beneficiaries
	housingRoomsRepository  repositories.HousingRooms
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindManyByHousingRoomIDPaginated(
	beneficiariesRepository repositories.Beneficiaries,
	housingRoomsRepository repositories.HousingRooms,
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) FindManyByHousingRoomIDPaginated {
	return &findManyByHousingRoomIDPaginatedImpl{
		beneficiariesRepository: beneficiariesRepository,
		housingRoomsRepository:  housingRoomsRepository,
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByHousingRoomIDPaginatedImpl) Execute(actor entities.User, roomID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
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

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.beneficiariesRepository.FindManyByRoomIDPaginated(room.ID, search, offset, limit)
}
