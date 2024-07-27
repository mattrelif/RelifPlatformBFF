package responses

import "relif/bff/entities"

type Address struct {
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	ZipCode      string `json:"zip_code"`
	District     string `json:"district"`
	City         string `json:"city"`
	Country      string `json:"country"`
}

func NewAddress(address entities.Address) Address {
	return Address{
		StreetName:   address.StreetName,
		StreetNumber: address.StreetNumber,
		ZipCode:      address.ZipCode,
		District:     address.District,
		City:         address.City,
		Country:      address.Country,
	}
}
