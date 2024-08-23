package organization_data_access_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User, targetOrganizationID string) (entities.OrganizationDataAccessRequest, error)
}

type createImpl struct {
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests
	organizationsRepository                  repositories.Organizations
}

func NewCreate(
	organizationDataAccessRequestsRepository repositories.OrganizationDataAccessRequests,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		organizationDataAccessRequestsRepository: organizationDataAccessRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, targetOrganizationID string) (entities.OrganizationDataAccessRequest, error) {
	targetOrganization, err := uc.organizationsRepository.FindOneByID(targetOrganizationID)

	if err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, actor.Organization); err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	request := entities.OrganizationDataAccessRequest{
		TargetOrganizationID: targetOrganization.ID,
		RequesterID:          actor.ID,
	}

	return uc.organizationDataAccessRequestsRepository.Create(request)
}
