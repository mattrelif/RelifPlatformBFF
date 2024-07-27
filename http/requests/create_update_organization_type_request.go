package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type CreateUpdateOrganizationTypeRequest struct {
	OrganizationID string `json:"organization_id"`
}

func (req *CreateUpdateOrganizationTypeRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.OrganizationID, validation.Required, is.MongoID),
	)
}

func (req *CreateUpdateOrganizationTypeRequest) ToEntity() entities.UpdateOrganizationTypeRequest {
	return entities.UpdateOrganizationTypeRequest{
		OrganizationID: req.OrganizationID,
	}
}
