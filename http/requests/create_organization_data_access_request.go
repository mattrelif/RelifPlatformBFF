package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type CreateOrganizationDataAccessRequest struct {
	TargetOrganizationID string `json:"target_organization_id"`
}

func (req *CreateOrganizationDataAccessRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.TargetOrganizationID, validation.Required, is.MongoID),
	)
}

func (req *CreateOrganizationDataAccessRequest) ToEntity() entities.OrganizationDataAccessRequest {
	return entities.OrganizationDataAccessRequest{
		TargetOrganizationID: req.TargetOrganizationID,
	}
}
