package organizations

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
)

type FindManyPaginated interface {
	Execute(offset, limit int64) (int64, []entities.Organization, error)
}

type findManyPaginatedImpl struct {
	organizationsRepository repositories.Organizations
}

func NewFindManyPaginated(organizationsRepository repositories.Organizations) FindManyPaginated {
	return &findManyPaginatedImpl{
		organizationsRepository: organizationsRepository,
	}
}

func (uc *findManyPaginatedImpl) Execute(offset, limit int64) (int64, []entities.Organization, error) {
	return uc.organizationsRepository.FindManyPaginated(offset, limit)
}
