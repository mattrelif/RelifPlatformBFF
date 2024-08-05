package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type RejectJoinOrganizationRequest struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectJoinOrganizationRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.RejectReason, validation.Required),
	)
}

func (req *RejectJoinOrganizationRequest) ToEntity() entities.JoinOrganizationRequest {
	return entities.JoinOrganizationRequest{
		RejectReason: req.RejectReason,
	}
}
