package services

import (
	"relif/bff/entities"
	"relif/bff/utils"
)

type Authorization interface {
	AuthorizePrivateActions(user entities.User) error
	AuthorizeAccessUserResource(userId string, user entities.User) error
	AuthorizeMutateUserData(userId string, user entities.User) error
	AuthorizeAccessOrganizationData(organizationId string, user entities.User) error
	AuthorizeAccessPrivateOrganizationData(organizationId string, user entities.User) error
	AuthorizeMutateOrganizationData(organizationId string, user entities.User) error
	AuthorizeCreateOrganization(user entities.User) error
	AuthorizeCreateOrganizationResource(user entities.User) error
	AuthorizeCreateAccessOrganizationDataRequest(user entities.User) error
	AuthorizeMutateAccessOrganizationDataRequestData(requestId string, user entities.User) error
	AuthorizeMutateOrganizationDataAccessGrantsData(grantId string, user entities.User) error
	AuthorizeAccessHousingData(housingId string, user entities.User) error
	AuthorizeMutateHousingData(housingId string, user entities.User) error
	AuthorizeCreateHousingResource(housingId string, user entities.User) error
	AuthorizeMutateJoinOrganizationInviteData(inviteId string, user entities.User) error
	AuthorizeCreateUpdateOrganizationTypeRequest(user entities.User) error
	AuthorizeCreateJoinOrganizationRequest(user entities.User) error
	AuthorizeMutateJoinOrganizationRequest(requestId string, user entities.User) error
	AuthorizeAccessHousingRoomData(roomId string, user entities.User) error
	AuthorizeMutateHousingRoomData(roomId string, user entities.User) error
	AuthorizeAccessBeneficiaryData(beneficiaryId string, user entities.User) error
	AuthorizeMutateBeneficiaryData(beneficiaryId string, user entities.User) error
	AuthorizeCreateBeneficiaryResource(beneficiaryId string, user entities.User) error
	AuthorizeAccessVoluntaryPersonData(voluntaryId string, user entities.User) error
	AuthorizeMutateVoluntaryPersonData(voluntaryId string, user entities.User) error
	AuthorizeAccessProductTypeData(typeId string, user entities.User) error
	AuthorizeMutateProductTypeData(typeId string, user entities.User) error
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

func (service *authorizationImpl) AuthorizeAccessUserResource(userId string, user entities.User) error {
	if user.ID != userId && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateUserData(userId string, user entities.User) error {
	target, err := service.usersService.FindOneById(userId)

	if err != nil {
		return err
	}

	if target.ID != user.ID && (target.OrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeAccessOrganizationData(organizationId string, user entities.User) error {
	accessGranted, err := service.organizationDataAccessGrantsService.ExistsByOrganizationIdAndTargetOrganizationId(user.OrganizationID, organizationId)

	if err != nil {
		return err
	}

	if user.OrganizationID != organizationId && !accessGranted {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeAccessPrivateOrganizationData(organizationId string, user entities.User) error {
	if (user.OrganizationID != organizationId && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateOrganizationData(organizationId string, user entities.User) error {
	if (user.OrganizationID != organizationId && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
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

func (service *authorizationImpl) AuthorizeCreateOrganizationResource(user entities.User) error {
	if (user.OrganizationID != "" && user.PlatformRole != utils.OrgAdminPlatformRole) && user.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeCreateAccessOrganizationDataRequest(user entities.User) error {
	organization, err := service.organizationsService.FindOneById(user.OrganizationID)

	if err != nil {
		return err
	}

	if organization.Type != utils.CoordinatorOrganizationType && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateAccessOrganizationDataRequestData(requestId string, user entities.User) error {
	request, err := service.organizationDataAccessRequestsService.FindOneById(requestId)

	if err != nil {
		return err
	}

	if request.RequesterOrganizationID != user.OrganizationID && user.PlatformRole != utils.OrgAdminPlatformRole {
		return utils.ErrUnauthorizedAction
	}

	return nil
}

func (service *authorizationImpl) AuthorizeMutateOrganizationDataAccessGrantsData(grantId string, user entities.User) error {
	grant, err := service.organizationDataAccessGrantsService.FindOneById(grantId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(grant.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessHousingData(housingId string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingId)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(housing.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateHousingData(housingId string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(housing.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeCreateHousingResource(housingId string, user entities.User) error {
	housing, err := service.housingsService.FindOneByID(housingId)

	if err != nil {
		return err
	}

	return service.AuthorizeCreateHousingResource(housing.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateJoinOrganizationInviteData(inviteId string, user entities.User) error {
	invite, err := service.joinOrganizationInvitesService.FindOneById(inviteId)

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

func (service *authorizationImpl) AuthorizeMutateJoinOrganizationRequest(requestId string, user entities.User) error {
	request, err := service.joinOrganizationRequestsService.FindOneById(requestId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(request.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessHousingRoomData(roomId string, user entities.User) error {
	room, err := service.housingRoomsService.FindOneById(roomId)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessHousingData(room.HousingID, user)
}

func (service *authorizationImpl) AuthorizeMutateHousingRoomData(roomId string, user entities.User) error {
	room, err := service.housingRoomsService.FindOneById(roomId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateHousingData(room.HousingID, user)
}

func (service *authorizationImpl) AuthorizeAccessBeneficiaryData(beneficiaryId string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneById(beneficiaryId)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(beneficiary.CurrentOrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateBeneficiaryData(beneficiaryId string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneById(beneficiaryId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(beneficiary.CurrentOrganizationID, user)
}

func (service *authorizationImpl) AuthorizeCreateBeneficiaryResource(beneficiaryId string, user entities.User) error {
	beneficiary, err := service.beneficiariesService.FindOneById(beneficiaryId)

	if err != nil {
		return err
	}

	if beneficiary.CurrentOrganizationID != user.OrganizationID {
		return utils.ErrUnauthorizedAction
	}

	return service.AuthorizeCreateOrganizationResource(user)
}

func (service *authorizationImpl) AuthorizeAccessVoluntaryPersonData(voluntaryId string, user entities.User) error {
	voluntary, err := service.voluntaryPeopleService.FindOneById(voluntaryId)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(voluntary.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateVoluntaryPersonData(voluntaryId string, user entities.User) error {
	voluntary, err := service.voluntaryPeopleService.FindOneById(voluntaryId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(voluntary.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeAccessProductTypeData(typeId string, user entities.User) error {
	productType, err := service.productTypesService.FindOneById(typeId)

	if err != nil {
		return err
	}

	return service.AuthorizeAccessOrganizationData(productType.OrganizationID, user)
}

func (service *authorizationImpl) AuthorizeMutateProductTypeData(typeId string, user entities.User) error {
	productType, err := service.productTypesService.FindOneById(typeId)

	if err != nil {
		return err
	}

	return service.AuthorizeMutateOrganizationData(productType.OrganizationID, user)
}
