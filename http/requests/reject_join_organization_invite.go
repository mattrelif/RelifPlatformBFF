package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type RejectJoinOrganizationInvite struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectJoinOrganizationInvite) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.RejectReason, validation.Required),
	)
}

func (req *RejectJoinOrganizationInvite) ToEntity() entities.JoinOrganizationInvite {
	return entities.JoinOrganizationInvite{
		RejectReason: req.RejectReason,
	}
}
