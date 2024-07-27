package models

import "relif/bff/entities"

type Address struct {
	StreetName   string `bson:"street_name"`
	StreetNumber string `bson:"street_number"`
	ZipCode      string `bson:"zip_code"`
	District     string `bson:"district"`
	City         string `bson:"city"`
	Country      string `bson:"country"`
}

func (address *Address) ToEntity() entities.Address {
	return entities.Address{
		StreetName:   address.StreetName,
		StreetNumber: address.StreetNumber,
		ZipCode:      address.ZipCode,
		District:     address.District,
		City:         address.City,
		Country:      address.Country,
	}
}
