package housing_rooms

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByHousingIDPaginated interface {
	Execute(actor entities.User, housingID string, offset, limit int64) (int64, []entities.HousingRoom, error)
}

type findManyByHousingIDPaginatedImpl struct {
	housingRoomsRepository  repositories.HousingRooms
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindManyByHousingIDPaginated(
	housingRoomsRepository repositories.HousingRooms,
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) FindManyByHousingIDPaginated {
	return &findManyByHousingIDPaginatedImpl{
		housingRoomsRepository:  housingRoomsRepository,
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByHousingIDPaginatedImpl) Execute(actor entities.User, housingID string, offset, limit int64) (int64, []entities.HousingRoom, error) {
	housing, err := uc.housingsRepository.FindOneByID(housingID)

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

	return uc.housingRoomsRepository.FindManyByHousingIDPaginated(housingID, offset, limit)
}
