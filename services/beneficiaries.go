package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
)

type Beneficiaries interface {
	Create(organizationId string, data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingId(housingId string) ([]entities.Beneficiary, error)
	FindManyByRoomId(roomId string) ([]entities.Beneficiary, error)
	FindManyByOrganizationId(organizationId string) ([]entities.Beneficiary, error)
	FindOneById(id string) (entities.Beneficiary, error)
	UpdateOneById(id string, data entities.Beneficiary) error
	InactivateOneById(id string) error
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

func (service *beneficiariesImpl) Create(organizationId string, data entities.Beneficiary) (entities.Beneficiary, error) {
	exists, err := service.ExistsByEmail(data.Email)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if exists {
		return entities.Beneficiary{}, utils.ErrBeneficiaryAlreadyExists
	}

	data.CurrentOrganizationID = organizationId

	return service.repository.Create(data)
}

func (service *beneficiariesImpl) FindManyByHousingId(housingId string) ([]entities.Beneficiary, error) {
	return service.repository.FindManyByHousingId(housingId)
}

func (service *beneficiariesImpl) FindManyByRoomId(roomId string) ([]entities.Beneficiary, error) {
	return service.repository.FindManyByRoomId(roomId)
}

func (service *beneficiariesImpl) FindManyByOrganizationId(roomId string) ([]entities.Beneficiary, error) {
	return service.repository.FindManyByOrganizationId(roomId)
}

func (service *beneficiariesImpl) FindOneById(id string) (entities.Beneficiary, error) {
	return service.repository.FindOneById(id)
}

func (service *beneficiariesImpl) UpdateOneById(id string, data entities.Beneficiary) error {
	return service.repository.UpdateOneById(id, data)
}

func (service *beneficiariesImpl) InactivateOneById(id string) error {
	data := entities.Beneficiary{
		Status: utils.InactiveStatus,
	}
	return service.repository.UpdateOneById(id, data)
}

func (service *beneficiariesImpl) ExistsByEmail(email string) (bool, error) {
	count, err := service.repository.CountByEmail(email)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (service *beneficiariesImpl) AuthorizeCreate(user entities.User) error {
	if user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}
