package organization_data_access_grants

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type DeleteOneByID interface {
	Execute(actor entities.User, grantID string) error
}

type deleteOneByIDImpl struct {
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants
	organizationsRepository                repositories.Organizations
}

func NewDeleteOneByIDImpl(
	organizationDataAccessGrantsRepository repositories.OrganizationDataAccessGrants,
	organizationsRepository repositories.Organizations,
) DeleteOneByID {
	return &deleteOneByIDImpl{
		organizationDataAccessGrantsRepository: organizationDataAccessGrantsRepository,
		organizationsRepository:                organizationsRepository,
	}
}

func (uc *deleteOneByIDImpl) Execute(actor entities.User, grantID string) error {
	grant, err := uc.organizationDataAccessGrantsRepository.FindOneByID(grantID)

	if err != nil {
		return err
	}

	targetOrganization, err := uc.organizationsRepository.FindOneByID(grant.TargetOrganizationID)

	if err != nil {
		return err
	}

	if err = guards.IsOrganizationAdmin(actor, targetOrganization); err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(grant.OrganizationID)

	if err != nil {
		return err
	}

	for index, id := range organization.AccessGrantedIDs {
		if id == targetOrganization.ID {
			organization.AccessGrantedIDs = append(organization.AccessGrantedIDs[:index], organization.AccessGrantedIDs[index+1:]...)
		}
	}

	if err = uc.organizationsRepository.UpdateOneByID(organization.ID, organization); err != nil {
		return err
	}

	if err = uc.organizationDataAccessGrantsRepository.DeleteOneByID(grant.ID); err != nil {
		return err
	}

	return nil
}
