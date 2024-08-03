package responses

import "relif/bff/entities"

type Address struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	ZipCode      string `json:"zip_code"`
	District     string `json:"district"`
	City         string `json:"city"`
	Country      string `json:"country"`
}

func NewAddress(address entities.Address) Address {
	return Address{
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		ZipCode:      address.ZipCode,
		District:     address.District,
		City:         address.City,
		Country:      address.Country,
	}
}
