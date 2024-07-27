package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type CreateJoinOrganizationInvite struct {
	OrganizationID string `json:"organization_id"`
}

func (req *CreateJoinOrganizationInvite) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.OrganizationID, validation.Required),
	)
}

func (req *CreateJoinOrganizationInvite) ToEntity() entities.JoinOrganizationInvite {
	return entities.JoinOrganizationInvite{
		OrganizationID: req.OrganizationID,
	}
}
