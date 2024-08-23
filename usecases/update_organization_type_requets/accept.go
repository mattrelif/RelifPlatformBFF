package update_organization_type_requets

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/utils"
	"time"
)

type Accept interface {
	Execute(actor entities.User, requestID string) error
}

type acceptImpl struct {
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests
	organizationsRepository                  repositories.Organizations
}

func NewAccept(
	updateOrganizationTypeRequestsRepository repositories.UpdateOrganizationTypeRequests,
	organizationsRepository repositories.Organizations,
) Accept {
	return &acceptImpl{
		updateOrganizationTypeRequestsRepository: updateOrganizationTypeRequestsRepository,
		organizationsRepository:                  organizationsRepository,
	}
}

func (uc *acceptImpl) Execute(actor entities.User, requestID string) error {
	if err := guards.IsSuperUser(actor); err != nil {
		return err
	}

	request, err := uc.updateOrganizationTypeRequestsRepository.FindOneByID(requestID)

	if err != nil {
		return err
	}

	organization, err := uc.organizationsRepository.FindOneByID(request.OrganizationID)

	if err != nil {
		return err
	}

	organization.Type = utils.CoordinatorOrganizationType

	if err = uc.organizationsRepository.UpdateOneByID(organization.ID, organization); err != nil {
		return err
	}

	request.Status = utils.AcceptedStatus
	request.AuditorID = actor.ID
	request.AcceptedAt = time.Now()

	if err = uc.updateOrganizationTypeRequestsRepository.UpdateOneByID(request.ID, request); err != nil {
		return err
	}

	return nil
}
