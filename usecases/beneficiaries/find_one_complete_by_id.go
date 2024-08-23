package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindOneCompleteByID interface {
	Execute(actor entities.User, id string) (entities.Beneficiary, error)
}

type findOneCompleteByIDImpl struct {
	beneficiariesRepository repositories.Beneficiaries
}

func NewFindOneCompleteByID(beneficiariesRepository repositories.Beneficiaries) FindOneCompleteByID {
	return &findOneCompleteByIDImpl{
		beneficiariesRepository: beneficiariesRepository,
	}
}

func (uc *findOneCompleteByIDImpl) Execute(actor entities.User, id string) (entities.Beneficiary, error) {
	beneficiary, err := uc.beneficiariesRepository.FindOneCompleteByID(id)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if err = guards.HasAccessToOrganizationData(actor, beneficiary.CurrentOrganization); err != nil {
		return entities.Beneficiary{}, err
	}

	return beneficiary, nil
}
