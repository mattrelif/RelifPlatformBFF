package organizations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type FindOneByID interface {
	Execute(id string) (entities.Organization, error)
}

type findOneByIDImpl struct {
	organizationsRepository repositories.Organizations
}

func NewFindOneByID(organizationsRepository repositories.Organizations) FindOneByID {
	return &findOneByIDImpl{
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findOneByIDImpl) Execute(id string) (entities.Organization, error) {
	return uc.organizationsRepository.FindOneByID(id)
}
