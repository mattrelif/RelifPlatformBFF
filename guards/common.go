package guards

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
)

func IsSuperUser(actor entities.User) error {
	if actor.PlatformRole != utils.RelifMemberPlatformRole {
		return utils.ErrForbiddenAction
	}

	return nil
}
