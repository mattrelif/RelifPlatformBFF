package update_organization_type_requets

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User) (entities.UpdateOrganizationTypeRequest, error)
}

type createImpl struct {
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests
	organizationsRepository                  repositories.Organizations
}

func NewCreate(
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests,
	organizationsRepository repositories.Organizations,
) Create {
	return &createImpl{
		updateOrganizationTypeRequestsRepository: updateOrganizationTypeRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User) (entities.UpdateOrganizationTypeRequest, error) {
	if err := guards.IsOrganizationAdmin(actor, actor.Organization); err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	request := entities.UpdateOrganizationTypeRequest{
		OrganizationID: actor.Organization.ID,
		CreatorID:      actor.ID,
	}

	return uc.updateOrganizationTypeRequestsRepository.Create(request)
}
