package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type CreateJoinOrganizationInvite struct {
	UserId string `json:"user_id"`
}

func (req *CreateJoinOrganizationInvite) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.UserId, validation.Required, is.MongoID),
	)
}

func (req *CreateJoinOrganizationInvite) ToEntity() entities.JoinOrganizationInvite {
	return entities.JoinOrganizationInvite{
		UserID: req.UserId,
	}
}
