package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	beneficiariesRepository repositories.Beneficiaries
	organizationsRepository repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		beneficiariesRepository: beneficiariesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.beneficiariesRepository.FindManyByOrganizationIDPaginated(organization.ID, search, offset, limit)
}
