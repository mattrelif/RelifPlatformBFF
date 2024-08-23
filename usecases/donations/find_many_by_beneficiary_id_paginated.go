package donations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByBeneficiaryIDPaginated interface {
	Execute(actor entities.User, beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error)
}

type findManyByBeneficiaryIDPaginated struct {
	donationsRepository     repositories.Donations
	beneficiariesRepository repositories.Beneficiaries
	organizationsRepository repositories.Organizations
}

func NewFindManyByBeneficiaryIDPaginated(
	donationsRepository repositories.Donations,
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) FindManyByBeneficiaryIDPaginated {
	return &findManyByBeneficiaryIDPaginated{
		donationsRepository:     donationsRepository,
		beneficiariesRepository: beneficiariesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByBeneficiaryIDPaginated) Execute(actor entities.User, beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error) {
	beneficiary, err := uc.beneficiariesRepository.FindOneByID(beneficiaryID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(beneficiary.CurrentOrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.HasAccessToOrganizationData(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.donationsRepository.FindManyByBeneficiaryIDPaginated(beneficiary.ID, offset, limit)
}
