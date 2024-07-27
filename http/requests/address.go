package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Address struct {
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	ZipCode      string `json:"zip_code"`
	District     string `json:"district"`
	City         string `json:"city"`
	Country      string `json:"country"`
}

func (req *Address) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.StreetName, validation.Required),
		validation.Field(&req.StreetNumber, validation.Required),
		validation.Field(&req.ZipCode, validation.Required),
		validation.Field(&req.District, validation.Required),
		validation.Field(&req.City, validation.Required),
		validation.Field(&req.Country, validation.Required),
	)
}
