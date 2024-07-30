package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type CreateProductType struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
}

func (req *CreateProductType) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.Brand, validation.Required),
		validation.Field(&req.Category, validation.Required),
	)
}

func (req *CreateProductType) ToEntity() entities.ProductType {
	return entities.ProductType{
		Name:        req.Name,
		Description: req.Description,
		Brand:       req.Brand,
		Category:    req.Category,
	}
}
