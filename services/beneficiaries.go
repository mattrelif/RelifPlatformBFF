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
	FindOneById(id string, user entities.User) (entities.Beneficiary, error)
	UpdateOneById(id string, data entities.Beneficiary) error
	InactivateOneById(id string) error
	ExistsByEmail(email string) (bool, error)
	AuthorizeCreate(user entities.User) error
	AuthorizeByHousingId(user entities.User, housingId string) error
	AuthorizeByRoomId(user entities.User, roomId string) error
	AuthorizeByOrganizationId(user entities.User, organizationId string) error
	AuthorizeExternalMutation(user entities.User, id string) error
}

type beneficiariesImpl struct {
	repository      repositories.Beneficiaries
	housingsService Housings
	roomsService    HousingRooms
	grantsService   OrganizationDataAccessGrants
}

func NewBeneficiaries(repository repositories.Beneficiaries, housingsService Housings, roomsService HousingRooms, grantsService OrganizationDataAccessGrants) Beneficiaries {
	return &beneficiariesImpl{
		repository:      repository,
		housingsService: housingsService,
		roomsService:    roomsService,
		grantsService:   grantsService,
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

func (service *beneficiariesImpl) FindOneById(id string, user entities.User) (entities.Beneficiary, error) {
	return service.authorizeFindOneById(id, user)
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

func (service *beneficiariesImpl) AuthorizeByHousingId(user entities.User, housingId string) error {
	if _, err := service.housingsService.FindOneByID(housingId, user); err != nil {
		return err
	}

	return nil
}

func (service *beneficiariesImpl) AuthorizeByRoomId(user entities.User, roomId string) error {
	if _, err := service.roomsService.FindOneById(roomId, user); err != nil {
		return err
	}

	return nil
}

func (service *beneficiariesImpl) AuthorizeByOrganizationId(user entities.User, organizationId string) error {
	accessGranted, err := service.grantsService.ExistsByOrganizationIdAndTargetOrganizationId(user.OrganizationID, organizationId)

	if err != nil {
		return err
	}

	if user.OrganizationID != organizationId && !accessGranted {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *beneficiariesImpl) AuthorizeExternalMutation(user entities.User, id string) error {
	beneficiary, err := service.repository.FindOneById(id)

	if err != nil {
		return err
	}

	accessGranted, err := service.grantsService.ExistsByOrganizationIdAndTargetOrganizationId(user.OrganizationID, beneficiary.CurrentOrganizationID)

	if err != nil {
		return err
	}

	if (user.OrganizationID != beneficiary.CurrentOrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && !accessGranted {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *beneficiariesImpl) authorizeFindOneById(id string, user entities.User) (entities.Beneficiary, error) {
	beneficiary, err := service.repository.FindOneById(id)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	if err = service.AuthorizeByOrganizationId(user, beneficiary.CurrentOrganizationID); err != nil {
		return entities.Beneficiary{}, err
	}

	return beneficiary, nil
}
