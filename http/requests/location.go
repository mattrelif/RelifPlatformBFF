package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/platform-bff/entities"
)

type Location struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

func (req *Location) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.MongoID),
		validation.Field(&req.Type, validation.Required),
	)
}

func (req *Location) ToEntity() entities.Location {
	return entities.Location{
		ID:   req.ID,
		Type: req.Type,
	}
}
