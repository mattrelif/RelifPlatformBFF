package update_organization_type_requets

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyByOrganizationIDPaginated interface {
	Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
}

type findManyByOrganizationIDPaginatedImpl struct {
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests
	organizationsRepository                  repositories.Organizations
}

func NewFindManyByOrganizationIDPaginated(
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests,
	organizationsRepository repositories.Organizations,
) FindManyByOrganizationIDPaginated {
	return &findManyByOrganizationIDPaginatedImpl{
		updateOrganizationTypeRequestsRepository: updateOrganizationTypeRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *findManyByOrganizationIDPaginatedImpl) Execute(actor entities.User, organizationID string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	organization, err := uc.organizationsRepository.FindOneByID(organizationID)

	if err != nil {
		return 0, nil, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return 0, nil, err
	}

	return uc.updateOrganizationTypeRequestsRepository.FindManyByOrganizationIDPaginated(organizationID, offset, limit)
}
