package models

import "relif/platform-bff/entities"

type Address struct {
	AddressLine1 string `bson:"address_line_1,omitempty"`
	AddressLine2 string `bson:"address_line_2,omitempty"`
	ZipCode      string `bson:"zip_code,omitempty"`
	District     string `bson:"district,omitempty"`
	City         string `bson:"city,omitempty"`
	Country      string `bson:"country,omitempty"`
}

func (address *Address) ToEntity() entities.Address {
	return entities.Address{
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		ZipCode:      address.ZipCode,
		District:     address.District,
		City:         address.City,
		Country:      address.Country,
	}
}

func NewAddress(entity entities.Address) Address {
	return Address{
		AddressLine1: entity.AddressLine1,
		AddressLine2: entity.AddressLine2,
		ZipCode:      entity.ZipCode,
		District:     entity.District,
		City:         entity.City,
		Country:      entity.Country,
	}
}
