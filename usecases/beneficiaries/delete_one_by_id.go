package beneficiaries

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, id string) error
}

type deleteOneByID struct {
	beneficiariesRepository repositories.Beneficiaries
	organizationsRepository repositories.Organizations
}

func NewDeleteOneByID(
	beneficiariesRepository repositories.Beneficiaries,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByID{
		beneficiariesRepository: beneficiariesRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *deleteOneByID) Execute(actor entities.User, id string) error {
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

	return uc.beneficiariesRepository.DeleteOneByID(beneficiary.ID)
}
