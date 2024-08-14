package services

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

type Authorization interface {
	AuthorizePrivateActions(user entities.User) error
	AuthorizeAccessUserResource(userID string, user entities.User) error
	AuthorizeMutateUserData(userID string, user entities.User) error
	AuthorizeAccessOrganizationData(organizationID string, user entities.User) error
	AuthorizeAccessPrivateOrganizationData(organizationID string, user entities.User) error
	AuthorizeMutateOrganizationData(organizationID string, user entities.User) error
	AuthorizeCreateOrganization(user entities.User) error
	AuthorizeCreateOrganizationResource(user entities.User, organizationID string) error
	AuthorizeCreateAccessOrganizationDataRequest(user entities.User) error
	AuthorizeMutateAccessOrganizationDataRequestData(requestID string, user entities.User) error
	AuthorizeMutateOrganizationDataAccessGrantsData(grantID string, user entities.User) error
	AuthorizeAccessHousingData(housingID string, user entities.User) error
	AuthorizeMutateHousingData(housingID string, user entities.User) error
	AuthorizeCreateHousingResource(housingID string, user entities.User) error
	AuthorizeMutateJoinOrganizationInviteData(inviteID string, user entities.User) error
	AuthorizeCreateUpdateOrganizationTypeRequest(user entities.User) error
	AuthorizeCreateJoinOrganizationRequest(user entities.User) error
	AuthorizeMutateJoinOrganizationRequest(requestID string, user entities.User) error
	AuthorizeAccessHousingRoomData(roomID string, user entities.User) error
	AuthorizeMutateHousingRoomData(roomID string, user entities.User) error
	AuthorizeAccessBeneficiaryData(beneficiaryID string, user entities.User) error
	AuthorizeMutateBeneficiaryData(beneficiaryID string, user entities.User) error
	AuthorizeCreateBeneficiaryResource(beneficiaryID string, user entities.User) error
	AuthorizeAccessVoluntaryPersonData(voluntaryID string, user entities.User) error
	AuthorizeMutateVoluntaryPersonData(voluntaryID string, user entities.User) error
	AuthorizeAccessProductTypeData(typeID string, user entities.User) error
	AuthorizeMutateProductTypeData(typeID string, user entities.User) error
	AuthorizeCreateProductTypeResource(productTypeID string, user entities.User) error
}

type authorizationImpl struct {
	usersService                          Users
	organizationsService                  Organizations
	housingsService                       Housings
	housingRoomsService                   HousingRooms
	beneficiariesService                  Beneficiaries
	organizationDataAccessRequestsService OrganizationDataAccessRequests
	organizationDataAccessGrantsService   OrganizationDataAccessGrants
	joinOrganizationInvitesService        JoinOrganizationInvites
	updateOrganizationTypeRequestsService UpdateOrganizationTypeRequests
	joinOrganizationRequestsService       JoinOrganizationRequests
	voluntaryPeopleService                VoluntaryPeople
	productTypesService                   ProductTypes
}

func NewAuthorization(
	usersService Users,
	organizationsService Organizations,
	housingsService Housings,
	housingRoomsService HousingRooms,
	beneficiariesService Beneficiaries,
	organizationDataAccessRequestsService OrganizationDataAccessRequests,
	organizationDataAccessGrantsService OrganizationDataAccessGrants,
	joinOrganizationInvitesService JoinOrganizationInvites,
	updateOrganizationTypeRequestsService UpdateOrganizationTypeRequests,
	joinOrganizationRequestsService JoinOrganizationRequests,
	voluntaryPeopleService VoluntaryPeople,
	productTypesService ProductTypes,
) Authorization {
	return &authorizationImpl{
		usersService:                          usersService,
		organizationsService:                  organizationsService,
		housingsService:                       housingsService,
		housingRoomsService:                   housingRoomsService,
		beneficiariesService:                  beneficiariesService,
		organizationDataAccessRequestsService: organizationDataAccessRequestsService,
		organizationDataAccessGrantsService:   organizationDataAccessGrantsService,
		joinOrganizationInvitesService:        joinOrganizationInvitesService,
		updateOrganizationTypeRequestsService: updateOrganizationTypeRequestsService,
		joinOrganizationRequestsService:       joinOrganizationRequestsService,
		voluntaryPeopleService:                voluntaryPeopleService,
		productTypesService:                   productTypesService,
	}
}

func (service *authorizationImpl) AuthorizePrivateActions(user entities.User) error {
	if user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeAccessUserResource(userID string, user entities.User) error {
	if user.ID != userID && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateUserData(userID string, user entities.User) error {
	target, err := service.usersService.FindOneByID(userID)

	if err != nil {
		return err
	}

	if target.ID != user.ID && (target.OrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeAccessOrganizationData(organizationID string, user entities.User) error {
	accessGranted, err := service.organizationDataAccessGrantsService.ExistsByOrganizationIDAndTargetOrganizationID(user.OrganizationID, organizationID)

	if err != nil {
		return err
	}

	if user.OrganizationID != organizationID && !accessGranted {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeAccessPrivateOrganizationData(organizationID string, user entities.User) error {
	if (user.OrganizationID != organizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateOrganizationData(organizationID string, user entities.User) error {
	if (user.OrganizationID != organizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateOrganization(user entities.User) error {
	if user.OrganizationID != "" && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateOrganizationResource(user entities.User, organizationID string) error {
	if (user.OrganizationID != organizationID && user.OrganizationID == "" && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateAccessOrganizationDataRequest(user entities.User) error {
	organization, err := service.organizationsService.FindOneByID(user.OrganizationID)

	if err != nil {
		return err
	}

	if organization.Type != utils.CoordinatorOrganizationType && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateAccessOrganizationDataRequestData(requestID string, user entities.User) error {
	request, err := service.organizationDataAccessRequestsService.FindOneByID(requestID)

	if err != nil {
		return err
	}

	if request.RequesterOrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateOrganizationDataAccessGrantsData(grantID string, user entities.User) error {
	grant, err := service.organizationDataAccessGrantsService.FindOneByID(grantID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(grant.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessHousingData(housingID string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingID)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(housing.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateHousingData(housingID string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(housing.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeCreateHousingResource(housingID string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingID)

	if err != nil {
		return err
	}

	if housing.OrganizationID != user.OrganizationID {
		return utils.ErrUnauthorizedAction
	}

	return service.AuthorizeCreateOrganizationResource(user, housing.OrganizationID)
}

func (service *authorizationImpl) AuthorizeMutateJoinOrganizationInviteData(inviteID string, user entities.User) error {
	invite, err := service.joinOrganizationInvitesService.FindOneByID(inviteID)

	if err != nil {
		return err
	}

	if (user.OrganizationID != "" && user.PlatformRole != utils.OrgMemberPlatformRole) && user.ID != invite.UserID {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateUpdateOrganizationTypeRequest(user entities.User) error {
	exists, err := service.updateOrganizationTypeRequestsService.ExistsPendingByOrganization(user.OrganizationID)

	if err != nil {
		return err
	}

	if user.PlatformRole != utils.OrgAdminPlatformRole && exists {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateJoinOrganizationRequest(user entities.User) error {
	if user.OrganizationID != "" {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateJoinOrganizationRequest(requestID string, user entities.User) error {
	request, err := service.joinOrganizationRequestsService.FindOneByID(requestID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(request.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessHousingRoomData(roomID string, user entities.User) error {
	room, err := service.housingRoomsService.FindOneByID(roomID)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessHousingData(room.HousingID, user)
}

func (service *authorizationImpl) AuthorizeMutateHousingRoomData(roomID string, user entities.User) error {
	room, err := service.housingRoomsService.FindOneByID(roomID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateHousingData(room.HousingID, user)
}

func (service *authorizationImpl) AuthorizeAccessBeneficiaryData(beneficiaryID string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(beneficiary.CurrentOrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateBeneficiaryData(beneficiaryID string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(beneficiary.CurrentOrganizationID, user)
}

func (service *authorizationImpl) AuthorizeCreateBeneficiaryResource(beneficiaryID string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneByID(beneficiaryID)

	if err != nil {
		return err
	}

	if beneficiary.CurrentOrganizationID != user.OrganizationID {
		return utils.ErrUnauthorizedAction
	}

	return service.AuthorizeCreateOrganizationResource(user, beneficiary.CurrentOrganizationID)
}

func (service *authorizationImpl) AuthorizeAccessVoluntaryPersonData(voluntaryID string, user entities.User) error {
	voluntary, err := service.voluntaryPeopleService.FindOneByID(voluntaryID)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(voluntary.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateVoluntaryPersonData(voluntaryID string, user entities.User) error {
	voluntary, err := service.voluntaryPeopleService.FindOneByID(voluntaryID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(voluntary.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessProductTypeData(typeID string, user entities.User) error {
	productType, err := service.productTypesService.FindOneByID(typeID)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(productType.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateProductTypeData(typeID string, user entities.User) error {
	productType, err := service.productTypesService.FindOneByID(typeID)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(productType.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeCreateProductTypeResource(productTypeID string, user entities.User) error {
	productType, err := service.productTypesService.FindOneByID(productTypeID)

	if err != nil {
		return err
	}

	return service.AuthorizeCreateOrganizationResource(user, productType.OrganizationID)
}
