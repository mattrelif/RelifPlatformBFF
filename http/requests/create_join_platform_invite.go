package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/platform-bff/entities"
)

type CreateJoinPlatformInvite struct {
	InvitedEmail string `json:"invited_email"`
}

func (req *CreateJoinPlatformInvite) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.InvitedEmail, validation.Required, is.Email),
	)
}

func (req *CreateJoinPlatformInvite) ToEntity() entities.JoinPlatformInvite {
	return entities.JoinPlatformInvite{
		InvitedEmail: req.InvitedEmail,
	}
}
