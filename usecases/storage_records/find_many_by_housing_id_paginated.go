package storage_records

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type FindManyByHousingIDPaginated interface {
	Execute(actor entities.User, housingID string, offset, limit int64) (int64, []entities.StorageRecord, error)
}

type findManyByHousingIDIDPaginatedImpl struct {
	storageRecordsRepository repositories.StorageRecords
	organizationsRepository  repositories.Organizations
	housingsRepository       repositories.Housings
}

func NewFindManyByHousingIDPaginated(
	storageRecordsRepository repositories.StorageRecords,
	organizationsRepository repositories.Organizations,
	housingsRepository repositories.Housings,
) FindManyByHousingIDPaginated {
	return &findManyByHousingIDIDPaginatedImpl{
		storageRecordsRepository: storageRecordsRepository,
		organizationsRepository:  organizationsRepository,
		housingsRepository:       housingsRepository,
	}
}

func (uc *findManyByHousingIDIDPaginatedImpl) Execute(actor entities.User, housingID string, offset, limit int64) (int64, []entities.StorageRecord, error) {
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

	location := entities.Location{
		ID:   housing.ID,
		Type: utils.HousingLocationType,
	}

	return uc.storageRecordsRepository.FindManyByLocationPaginated(location, offset, limit)
}
