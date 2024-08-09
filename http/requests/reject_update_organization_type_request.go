package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type RejectUpdateOrganizationTypeRequest struct {
	RejectReason string `json:"reject_reason"`
}

func (req *RejectUpdateOrganizationTypeRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.RejectReason, validation.Required),
	)
}

func (req *RejectUpdateOrganizationTypeRequest) ToEntity() entities.UpdateOrganizationTypeRequest {
	return entities.UpdateOrganizationTypeRequest{
		RejectReason: req.RejectReason,
	}
}
