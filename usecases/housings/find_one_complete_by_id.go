package housings

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type FindOneCompleteByID interface {
	Execute(actor entities.User, id string) (entities.Housing, error)
}

type findOneCompleteByIDImpl struct {
	housingsRepository      repositories.Housings
	organizationsRepository repositories.Organizations
}

func NewFindOneCompleteByID(
	housingsRepository repositories.Housings,
	organizationsRepository repositories.Organizations,
) FindOneCompleteByID {
	return &findOneCompleteByIDImpl{
		housingsRepository:      housingsRepository,
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findOneCompleteByIDImpl) Execute(actor entities.User, id string) (entities.Housing, error) {
	housing, err := uc.housingsRepository.FindOneCompleteByID(id)

	if err != nil {
		return entities.Housing{}, err
	}

	organization, err := uc.organizationsRepository.FindOneByID(housing.OrganizationID)

	if err != nil {
		return entities.Housing{}, err
	}

	if err = guards.IsOrganizationAdmin(actor, organization); err != nil {
		return entities.Housing{}, err
	}

	return housing, nil
}
