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
		"Has RELIF_MEMBER platform role - allowed": {
			user: entities.User{
				PlatformRole: utils.RelifMemberPlatformRole,
			},
			expectedError: nil,
		},
		"Has NO_ORG platform role - not allowed": {
			user: entities.User{
				PlatformRole: utils.NoOrgPlatformRole,
			},
			expectedError: utils.ErrForbiddenAction,
		},
		"Has ORG_MEMBER platform role - not allowed": {
			user: entities.User{
				PlatformRole: utils.OrgMemberPlatformRole,
			},
			expectedError: utils.ErrForbiddenAction,
		},
		"Has ORG_ADMIN platform role - not allowed": {
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
