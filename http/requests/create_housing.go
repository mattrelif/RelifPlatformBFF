package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/bff/entities"
)

type CreateHousing struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func (req *CreateHousing) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Address, validation.By(func(value interface{}) error {
			if address, ok := value.(Address); ok {
				return address.Validate()
			}

			return nil
		})),
	)
}

func (req *CreateHousing) ToEntity() entities.Housing {
	return entities.Housing{
		Name: req.Name,
		Address: entities.Address{
			StreetNumber: req.Address.StreetNumber,
			StreetName:   req.Address.StreetName,
			City:         req.Address.City,
			ZipCode:      req.Address.ZipCode,
			District:     req.Address.District,
			Country:      req.Address.Country,
		},
	}
}
