package housings

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
)

type Create interface {
	Execute(actor entities.User, data entities.Housing) (entities.Housing, error)
}

type createImpl struct {
	housingsRepository repositories.Housings
}

func NewCreate(housingsRepository repositories.Housings) Create {
	return &createImpl{
		housingsRepository: housingsRepository,
	}
}

func (uc *createImpl) Execute(actor entities.User, data entities.Housing) (entities.Housing, error) {
	if err := guards.IsOrganizationAdmin(actor, actor.Organization); err != nil {
		return entities.Housing{}, err
	}

	data.OrganizationID = actor.OrganizationID

	housing, err := uc.housingsRepository.Create(data)

	if err != nil {
		return entities.Housing{}, err
	}

	return housing, nil
}
