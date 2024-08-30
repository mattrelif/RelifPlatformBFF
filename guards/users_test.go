package guards

import (
	"github.com/stretchr/testify/assert"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"testing"
)

func Test_CanAccessPlatform(t *testing.T) {
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
		"Is active, is not a member of an organization and has the NO_ORG platform role - allowed": {
			user: entities.User{
				PlatformRole: utils.NoOrgPlatformRole,
				Status:       utils.ActiveStatus,
				Organization: entities.Organization{
					ID: "",
				},
			},
			expectedError: nil,
		},
		"Is active and member of an active organization - allowed": {
			user: entities.User{
				PlatformRole: utils.OrgMemberPlatformRole,
				Status:       utils.ActiveStatus,
				Organization: entities.Organization{
					ID:     "1",
					Status: utils.ActiveStatus,
				},
			},
		},
		"Is inactive - not allowed": {
			user: entities.User{
				Status: utils.InactiveStatus,
			},
			expectedError: utils.ErrInactiveUser,
		},
		"Is active and member of an inactive organization - not allowed": {
			user: entities.User{
				PlatformRole: utils.OrgMemberPlatformRole,
				Status:       utils.ActiveStatus,
				Organization: entities.Organization{
					ID:     "1",
					Status: utils.InactiveStatus,
				},
			},
			expectedError: utils.ErrMemberOfInactiveOrganization,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := CanAccessPlatform(test.user)
			assert.Equal(t, test.expectedError, err)
		})
	}
}

func Test_IsUser(t *testing.T) {
	tests := map[string]struct {
		actor         entities.User
		target        entities.User
		expectedError error
	}{
		"Actor has RELIF_MEMBER platform role - allowed": {
			actor: entities.User{
				PlatformRole: utils.RelifMemberPlatformRole,
			},
			target:        entities.User{},
			expectedError: nil,
		},
		"Actor has the same id as target - allowed": {
			actor: entities.User{
				ID:           "foo",
				PlatformRole: utils.NoOrgPlatformRole,
			},
			target: entities.User{
				ID: "foo",
			},
			expectedError: nil,
		},
		"Actor does not have the same id as target - not allowed": {
			actor: entities.User{
				ID:           "foo",
				PlatformRole: utils.NoOrgPlatformRole,
			},
			target: entities.User{
				ID: "bar",
			},
			expectedError: utils.ErrForbiddenAction,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := IsUser(test.actor, test.target)
			assert.Equal(t, test.expectedError, err)
		})
	}
}
