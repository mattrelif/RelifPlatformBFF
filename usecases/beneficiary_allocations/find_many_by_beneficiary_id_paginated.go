package beneficiary_allocations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByBeneficiaryIDPaginated interface {
	Execute(actor entities.User, beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
}

type findManyByBeneficiaryIDPaginatedImpl struct {
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations
	beneficiariesRepository          repositories.Beneficiaries
	organizationsRepository          repositories.Organizations
}

func NewFindManyByBeneficiaryIDPaginated(
	beneficiaryAllocationsRepository repositories.BeneficiaryAllocations,
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) FindManyByBeneficiaryIDPaginated {
	return &findManyByBeneficiaryIDPaginatedImpl{
		beneficiaryAllocationsRepository: beneficiaryAllocationsRepository,
		beneficiariesRepository:          beneficiariesRepository,
		organizationsRepository:          organizationsRepository,
	}
}

func (uc *findManyByBeneficiaryIDPaginatedImpl) Execute(actor entities.User, beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
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

	return uc.beneficiaryAllocationsRepository.FindManyByBeneficiaryIDPaginated(beneficiary.ID, offset, limit)
}
