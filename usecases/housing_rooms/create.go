package housing_rooms

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type CreateMany interface {
	Execute(actor entities.User, housingID string, data []entities.HousingRoom) (int64, []entities.HousingRoom, error)
}

type createManyImpl struct {
	housingRoomsRepository repositories.HousingRooms
	housingsRepository     repositories.Housings
	organizationRepository repositories.Organizations
}

func NewCreateHousingRoom(
	housingRoomsRepository repositories.HousingRooms,
	housingsRepository repositories.Housings,
	organizationRepository repositories.Organizations,
) CreateMany {
	return &createManyImpl{
		housingRoomsRepository: housingRoomsRepository,
		housingsRepository:     housingsRepository,
		organizationRepository: organizationRepository,
	}
}

func (uc *createManyImpl) Execute(actor entities.User, housingID string, data []entities.HousingRoom) (int64, []entities.HousingRoom, error) {
	housing, err := uc.housingsRepository.FindOneByID(housingID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	rooms, err := uc.housingRoomsRepository.CreateMany(data, housingID)

	if err != nil {
		return 0, nil, err
	}

	return int64(len(rooms)), rooms, nil
}
