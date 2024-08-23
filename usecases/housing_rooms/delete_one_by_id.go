package housing_rooms

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, id string) error
}

type deleteOneByIDImpl struct {
	housingRoomsRepository  repositories.HousingRooms
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewDeleteOneByID(
	housingRoomsRepository repositories.HousingRooms,
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByIDImpl{
		housingRoomsRepository:  housingRoomsRepository,
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *deleteOneByIDImpl) Execute(actor entities.User, id string) error {
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

	return uc.housingRoomsRepository.DeleteOneByID(room.ID)
}
