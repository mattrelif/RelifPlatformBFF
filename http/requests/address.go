package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	ZipCode      string `json:"zip_code"`
	District     string `json:"district"`
	City         string `json:"city"`
	Country      string `json:"country"`
}

func (req *Address) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.AddressLine1, validation.Required),
		validation.Field(&req.ZipCode, validation.Required),
		validation.Field(&req.District, validation.Required),
		validation.Field(&req.City, validation.Required),
		validation.Field(&req.Country, validation.Required),
	)
}

func (req *Address) ToEntity() entities.Address {
	return entities.Address{
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		ZipCode:      req.ZipCode,
		District:     req.District,
		City:         req.City,
		Country:      req.Country,
	}
}
