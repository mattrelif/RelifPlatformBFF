package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type RejectOrganizationDataAccessRequest struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectOrganizationDataAccessRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.RejectReason, validation.Required),
	)
}

func (req *RejectOrganizationDataAccessRequest) ToEntity() entities.OrganizationDataAccessRequest {
	return entities.OrganizationDataAccessRequest{
		RejectReason: req.RejectReason,
	}
}
