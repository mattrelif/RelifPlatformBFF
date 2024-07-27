package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type CreateJoinOrganizationRequest struct {
	OrganizationID string `json:"organization_id"`
}

func (req *CreateJoinOrganizationRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.OrganizationID, validation.Required),
	)
}

func (req *CreateJoinOrganizationRequest) ToEntity() entities.JoinOrganizationRequest {
	return entities.JoinOrganizationRequest{
		OrganizationID: req.OrganizationID,
	}
}
