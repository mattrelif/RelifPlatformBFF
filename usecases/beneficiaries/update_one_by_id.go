package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type UpdateOneByID interface {
	Execute(actor entities.User, id string, data entities.Beneficiary) error
}

type updateOneByIDImpl struct {
	beneficiariesRepository repositories.Beneficiaries
	organizationsRepository repositories.Organizations
}

func NewUpdateOneByID(
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) UpdateOneByID {
	return &updateOneByIDImpl{
		beneficiariesRepository: beneficiariesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *updateOneByIDImpl) Execute(actor entities.User, id string, data entities.Beneficiary) error {
	beneficiary, err := uc.beneficiariesRepository.FindOneByID(id)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(beneficiary.CurrentOrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return err
	}

	return uc.beneficiariesRepository.UpdateOneByID(beneficiary.ID, data)
}
