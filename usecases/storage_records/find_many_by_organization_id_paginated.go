package storage_records

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.StorageRecord, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	storageRecordsRepository repositories.StorageRecords
	organizationsRepository  repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	storageRecordsRepository repositories.StorageRecords,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		storageRecordsRepository: storageRecordsRepository,
		organizationsRepository:  organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.StorageRecord, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	location := entities.Location{
		ID:   organization.ID,
		Type: utils.OrganizationLocationType,
	}

	return uc.storageRecordsRepository.FindManyByLocationPaginated(location, offset, limit)
}
