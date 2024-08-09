package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type UpdateHousing struct {
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

func (req *UpdateHousing) Validate() error {
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

func (req *UpdateHousing) ToEntity() entities.Housing {
	return entities.Housing{
		Name:    req.Name,
		Address: req.Address.ToEntity(),
	}
}
