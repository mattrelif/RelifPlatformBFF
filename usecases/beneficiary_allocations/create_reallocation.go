package beneficiary_allocations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type CreateReallocation interface {
	Execute(actor entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
}

type createReallocationImpl struct {
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations
	beneficiariesRepository          repositories.Beneficiaries
	organizationsRepository          repositories.Organizations
	housingsRepository               repositories.Housings
	housingRoomsRepository           repositories.HousingRooms
}

func NewCreateReallocation(
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations,
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
	housingsRepository repositories.Housings,
	housingRoomsRepository repositories.HousingRooms,
) CreateReallocation {
	return &createReallocationImpl{
		beneficiaryAllocationsRepository: beneficiaryAllocationsRepository,
		beneficiariesRepository:          beneficiariesRepository,
		organizationsRepository:          organizationsRepository,
		housingsRepository:               housingsRepository,
		housingRoomsRepository:           housingRoomsRepository,
	}
}

func (uc *createReallocationImpl) Execute(actor entities.User, beneficiaryID string, data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	beneficiary, err := uc.beneficiariesRepository.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(beneficiary.CurrentOrganizationID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	data.OldRoomID = beneficiary.CurrentRoomID
	data.OldHousingID = beneficiary.CurrentHousingID

	housing, err := uc.housingsRepository.FindOneByID(data.HousingID)

	if err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	beneficiary.CurrentHousingID = housing.ID

	if data.RoomID != "" {
		var room entities.HousingRoom

		room, err = uc.housingRoomsRepository.FindOneByID(data.RoomID)

		if err != nil {
			return entities.BeneficiaryAllocation{}, err
		}

		beneficiary.CurrentRoomID = room.ID
	}

	if err = uc.beneficiariesRepository.UpdateOneByID(beneficiary.ID, beneficiary); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	data.AuditorID = actor.ID
	data.BeneficiaryID = beneficiary.ID
	data.Type = utils.ReallocationType

	return uc.beneficiaryAllocationsRepository.Create(data)
}
