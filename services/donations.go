package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type Donations interface {
	Create(beneficiaryID string, data entities.Donation) (entities.Donation, error)
	FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error)
}

type donationsImpl struct {
	repository                repositories.Donations
	productsInStoragesService ProductsInStorages
	beneficiariesService      Beneficiaries
}

func NewDonations(repository repositories.Donations, productsInStoragesService ProductsInStorages, beneficiariesService Beneficiaries) Donations {
	return &donationsImpl{
		repository:                repository,
		productsInStoragesService: productsInStoragesService,
		beneficiariesService:      beneficiariesService,
	}
}

func (service *donationsImpl) Create(beneficiaryID string, data entities.Donation) (entities.Donation, error) {
	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return entities.Donation{}, err
	}

	ids, err := service.productsInStoragesService.FindManyIDsByLocation(data.From, data.Quantity)

	if err != nil {
		return entities.Donation{}, err
	}

	if err = service.productsInStoragesService.DeleteManyByIDs(ids); err != nil {
		return entities.Donation{}, err
	}

	data.OrganizationID = beneficiary.CurrentOrganizationID

	donation, err := service.repository.Create(data)

	if err != nil {
		return entities.Donation{}, err
	}

	return donation, nil
}

func (service *donationsImpl) FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error) {
	return service.repository.FindManyByBeneficiaryID(beneficiaryID, offset, limit)
}
