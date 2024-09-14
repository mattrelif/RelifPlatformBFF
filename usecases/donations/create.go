package donations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, beneficiaryID string, data entities.Donation) (entities.Donation, error)
}

type createImpl struct {
	donationsRepository      repositories.Donations
	beneficiariesRepository  repositories.Beneficiaries
	storageRecordsRepository repositories.StorageRecords
	organizationsRepository  repositories.Organizations
	productTypesRepository   repositories.ProductTypes
}

func NewCreate(
	donationsRepository repositories.Donations,
	beneficiariesRepository repositories.Beneficiaries,
	storageRecordsRepository repositories.StorageRecords,
	organizationsRepository repositories.Organizations,
	productTypesRepository repositories.ProductTypes,
) Create {
	return &createImpl{
		donationsRepository:      donationsRepository,
		beneficiariesRepository:  beneficiariesRepository,
		storageRecordsRepository: storageRecordsRepository,
		organizationsRepository:  organizationsRepository,
		productTypesRepository:   productTypesRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, beneficiaryID string, data entities.Donation) (entities.Donation, error) {
	beneficiary, err := uc.beneficiariesRepository.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.Donation{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(beneficiary.CurrentOrganizationID)

	if err != nil {
		return entities.Donation{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.Donation{}, err
	}

	productType, err := uc.productTypesRepository.FindOneByID(data.ProductTypeID)

	if err != nil {
		return entities.Donation{}, err
	}

	record, err := uc.storageRecordsRepository.FindOneByProductTypeIDAndLocation(productType.ID, data.From)

	if err != nil {
		return entities.Donation{}, err
	}

	if record.ID != "" {
		if err = uc.storageRecordsRepository.DecreaseQuantityOfOneByID(record.ID, data.Quantity); err != nil {
			return entities.Donation{}, err
		}
	} else {
		return entities.Donation{}, utils.ErrStorageRecordNotFound
	}

	data.BeneficiaryID = beneficiary.ID
	data.OrganizationID = organization.ID

	return uc.donationsRepository.Create(data)
}
