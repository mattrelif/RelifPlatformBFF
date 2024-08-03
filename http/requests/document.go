package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type Document struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (req *Document) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Type, validation.Required),
		validation.Field(&req.Value, validation.Required),
	)
}

func (req *Document) ToEntity() entities.Document {
	return entities.Document{
		Type:  req.Type,
		Value: req.Value,
	}
}
