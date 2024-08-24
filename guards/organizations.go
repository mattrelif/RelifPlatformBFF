package guards

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

func AuthorizeCreateOrganization(actor entities.User) error {
	if err := IsSuperUser(actor); err != nil {
		if actor.Organization.ID != "" && actor.PlatformRole != utils.NoOrgPlatformRole {
			return utils.ErrForbiddenAction
		}
	}

	return nil
}

func IsOrganizationAdmin(actor entities.User, organization entities.Organization) error {
	if err := IsSuperUser(actor); err != nil {
		if actor.Organization.ID != organization.ID && actor.PlatformRole != utils.OrgAdminPlatformRole {
			return utils.ErrForbiddenAction
		}
	}

	return nil
}

func organizationHasGrant(actor entities.User, target entities.Organization) bool {
	for _, id := range actor.Organization.AccessGrantedIDs {
		if id == target.ID {
			return true
		}
	}

	return false
}

func HasAccessToOrganizationData(actor entities.User, target entities.Organization) error {
	if err := IsSuperUser(actor); err != nil {
		if actor.Organization.ID != target.ID && !organizationHasGrant(actor, target) {
			return utils.ErrForbiddenAction
		}
	}

	return nil
}
