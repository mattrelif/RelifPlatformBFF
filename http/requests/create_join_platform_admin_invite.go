package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/platform-bff/entities"
)

type CreateJoinPlatformAdminInvite struct {
	InvitedEmail string `json:"invited_email"`
}

func (req *CreateJoinPlatformAdminInvite) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.InvitedEmail, validation.Required, is.Email),
	)
}

func (req *CreateJoinPlatformAdminInvite) ToEntity() entities.JoinPlatformAdminInvite {
	return entities.JoinPlatformAdminInvite{
		InvitedEmail: req.InvitedEmail,
	}
}
