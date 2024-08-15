package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Donations interface {
	Create(beneficiaryID string, data entities.Donation) (entities.Donation, error)
	FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error)
}

type donationsImpl struct {
	repository            repositories.Donations
	storageRecordsService StorageRecords
	beneficiariesService  Beneficiaries
}

func NewDonations(repository repositories.Donations, storageRecordsService StorageRecords, beneficiariesService Beneficiaries) Donations {
	return &donationsImpl{
		repository:            repository,
		storageRecordsService: storageRecordsService,
		beneficiariesService:  beneficiariesService,
	}
}

func (service *donationsImpl) Create(beneficiaryID string, data entities.Donation) (entities.Donation, error) {
	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.Donation{}, err
	}

	record, err := service.storageRecordsService.FindOneByProductTypeIDAndLocation(data.ProductTypeID, data.From)

	if err != nil {
		return entities.Donation{}, err
	}

	if record.ID != "" {
		record.Quantity -= data.Quantity

		if err = service.storageRecordsService.UpdateOneByID(record.ID, record); err != nil {
			return entities.Donation{}, err
		}
	} else {
		return entities.Donation{}, utils.ErrStorageRecordNotFound
	}

	data.OrganizationID = beneficiary.CurrentOrganizationID

	return service.repository.Create(data)
}

func (service *donationsImpl) FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error) {
	return service.repository.FindManyByBeneficiaryID(beneficiaryID, offset, limit)
}
