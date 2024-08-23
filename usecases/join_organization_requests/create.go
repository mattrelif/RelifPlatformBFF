package join_organization_requests

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User, organizationID string) (entities.JoinOrganizationRequest, error)
}

type createImpl struct {
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests
	organizationsRepository            repositories.Organizations
}

func NewCreate(
	joinOrganizationRequestsRepository repositories.JoinOrganizationRequests,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		joinOrganizationRequestsRepository: joinOrganizationRequestsRepository,
		organizationsRepository:            organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, organizationID string) (entities.JoinOrganizationRequest, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return entities.JoinOrganizationRequest{}, err
	}

	request := entities.JoinOrganizationRequest{
		UserID:         actor.ID,
		OrganizationID: organization.ID,
	}

	return uc.joinOrganizationRequestsRepository.Create(request)
}
