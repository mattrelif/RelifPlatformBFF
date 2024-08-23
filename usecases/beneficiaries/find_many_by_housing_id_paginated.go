package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByHousingIDPaginated interface {
	Execute(actor entities.User, housingID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
}

type findManyByHousingIDPaginatedImpl struct {
	beneficiariesRepository repositories.Beneficiaries
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindManyByHousingIDPaginated(
	beneficiariesRepository repositories.Beneficiaries,
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) FindManyByHousingIDPaginated {
	return &findManyByHousingIDPaginatedImpl{
		beneficiariesRepository: beneficiariesRepository,
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyByHousingIDPaginatedImpl) Execute(actor entities.User, housingID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
	housing, err := uc.housingsRepository.FindOneByID(housingID)

	if err != nil {
		return 0, nil, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.beneficiariesRepository.FindManyByHousingIDPaginated(housing.ID, search, offset, limit)
}
