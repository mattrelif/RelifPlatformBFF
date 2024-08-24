package guards

import (
	"github.com/stretchr/testify/assert"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"testing"
)

func Test_IsSuperUser(t *testing.T) {
	tests := map[string]struct {
		user          entities.User
		expectedError error
	}{
		"Happy path": {
			user: entities.User{
				PlatformRole: utils.RelifMemberPlatformRole,
			},
			expectedError: nil,
		},
		"NO_ORG role - Invalid": {
			user: entities.User{
				PlatformRole: utils.NoOrgPlatformRole,
			},
			expectedError: utils.ErrForbiddenAction,
		},
		"ORG_MEMBER role - Invalid": {
			user: entities.User{
				PlatformRole: utils.OrgMemberPlatformRole,
			},
			expectedError: utils.ErrForbiddenAction,
		},
		"ORG_ADMIN role - Invalid": {
			user: entities.User{
				PlatformRole: utils.OrgAdminPlatformRole,
			},
			expectedError: utils.ErrForbiddenAction,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := IsSuperUser(test.user)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
