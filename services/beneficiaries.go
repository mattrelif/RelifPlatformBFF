package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
)

type Beneficiaries interface {
	Create(organizationID string, data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingID(housingID, search string, limit, offset int64) (int64, []entities.Beneficiary, error)
	FindManyByRoomID(roomID, search string, limit, offset int64) (int64, []entities.Beneficiary, error)
	FindManyByOrganizationID(organizationID, search string, limit, offset int64) (int64, []entities.Beneficiary, error)
	FindOneByID(id string) (entities.Beneficiary, error)
	UpdateOneByID(id string, data entities.Beneficiary) error
	InactivateOneByID(id string) error
	ExistsByEmail(email string) (bool, error)
}

type beneficiariesImpl struct {
	repository repositories.Beneficiaries
}

func NewBeneficiaries(repository repositories.Beneficiaries) Beneficiaries {
	return &beneficiariesImpl{
		repository: repository,
	}
}

func (service *beneficiariesImpl) Create(organizationID string, data entities.Beneficiary) (entities.Beneficiary, error) {
	exists, err := service.ExistsByEmail(data.Email)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if exists {
		return entities.Beneficiary{}, utils.ErrBeneficiaryAlreadyExists
	}

	data.CurrentOrganizationID = organizationID

	return service.repository.Create(data)
}

func (service *beneficiariesImpl) FindManyByHousingID(housingID, search string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	return service.repository.FindManyByHousingID(housingID, search, limit, offset)
}

func (service *beneficiariesImpl) FindManyByRoomID(roomID, search string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	return service.repository.FindManyByRoomID(roomID, search, limit, offset)
}

func (service *beneficiariesImpl) FindManyByOrganizationID(roomID, search string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	return service.repository.FindManyByOrganizationID(roomID, search, limit, offset)
}

func (service *beneficiariesImpl) FindOneByID(id string) (entities.Beneficiary, error) {
	return service.repository.FindOneByID(id)
}

func (service *beneficiariesImpl) UpdateOneByID(id string, data entities.Beneficiary) error {
	return service.repository.UpdateOneByID(id, data)
}

func (service *beneficiariesImpl) InactivateOneByID(id string) error {
	data := entities.Beneficiary{
		Status: utils.InactiveStatus,
	}
	return service.repository.UpdateOneByID(id, data)
}

func (service *beneficiariesImpl) ExistsByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
