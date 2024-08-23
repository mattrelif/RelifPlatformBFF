package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Create interface {
	Execute(actor entities.User, organizationID string, data entities.Beneficiary) (entities.Beneficiary, error)
}

type createImpl struct {
	beneficiariesRepository repositories.Beneficiaries
	organizationsRepository repositories.Organizations
}

func NewCreate(
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		beneficiariesRepository: beneficiariesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, organizationID string, data entities.Beneficiary) (entities.Beneficiary, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.Beneficiary{}, err
	}

	count, err := uc.beneficiariesRepository.CountByEmail(data.Email)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if count > 0 {
		return entities.Beneficiary{}, utils.ErrBeneficiaryAlreadyExists
	}
	
	data.CurrentOrganizationID = organization.ID

	return uc.beneficiariesRepository.Create(data)
}
