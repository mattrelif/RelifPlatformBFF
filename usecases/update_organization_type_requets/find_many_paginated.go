package update_organization_type_requets

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindManyPaginated interface {
	Execute(actor entities.User, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
}

type findManyPaginatedImpl struct {
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests
}

func NewFindManyPaginated(updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests) FindManyPaginated {
	return &findManyPaginatedImpl{
		updateOrganizationTypeRequestsRepository: updateOrganizationTypeRequestsRepository,
	}
}

func (uc *findManyPaginatedImpl) Execute(actor entities.User, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	if err := guards.IsSuperUser(actor); err != nil {
		return 0, nil, err
	}

	return uc.updateOrganizationTypeRequestsRepository.FindManyPaginated(offset, limit)
}
