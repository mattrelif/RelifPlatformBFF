package housing_rooms

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindOneCompleteByID interface {
	Execute(actor entities.User, id string) (entities.HousingRoom, error)
}

type findOneCompleteByIDImpl struct {
	housingRoomsRepository  repositories.HousingRooms
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindOneCompleteByID(
	housingRoomsRepository repositories.HousingRooms,
	organizationsRepository repositories.Organizations,
) FindOneCompleteByID {
	return &findOneCompleteByIDImpl{
		housingRoomsRepository:  housingRoomsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findOneCompleteByIDImpl) Execute(actor entities.User, id string) (entities.HousingRoom, error) {
	room, err := uc.housingRoomsRepository.FindOneCompleteByID(id)

	if err != nil {
		return entities.HousingRoom{}, err
	}

	housing, err := uc.housingsRepository.FindOneByID(room.HousingID)

	if err != nil {
		return entities.HousingRoom{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return entities.HousingRoom{}, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return entities.HousingRoom{}, err
	}

	return room, nil
}
