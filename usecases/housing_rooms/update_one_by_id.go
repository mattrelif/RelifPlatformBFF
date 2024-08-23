package housing_rooms

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.HousingRoom) error
}

type updateOneByIDImpl struct {
	housingRoomsRepository  repositories.HousingRooms
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewUpdateOneByID(
	housingRoomsRepository repositories.HousingRooms,
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) UpdateOneByID {
	return &updateOneByIDImpl{
		housingRoomsRepository:  housingRoomsRepository,
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.HousingRoom) error {
	room, err := uc.housingRoomsRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	housing, err := uc.housingsRepository.FindOneByID(room.HousingID)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	return uc.housingRoomsRepository.UpdateOneByID(room.ID, data)
}
