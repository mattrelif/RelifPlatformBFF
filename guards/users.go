package guards

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

func CanAccessPlatform(actor entities.User) error {
	if actor.Status == utils.InactiveStatus {
		return utils.ErrInactiveUser
	}

	if actor.Organization.ID != "" {
		if actor.Organization.Status == utils.InactiveStatus {
			return utils.ErrMemberOfInactiveOrganization
		}
	}

	return nil
}

func IsUser(actor, target entities.User) error {
	if err := IsSuperUser(actor); err != nil {
		if actor.ID != target.ID {
			return utils.ErrForbiddenAction
		}
	}

	return nil
}
