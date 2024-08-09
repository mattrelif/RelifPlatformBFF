package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type CreateOrganization struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     Address `json:"address"`
}

func (req *CreateOrganization) Validate() error {
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

func (req *CreateOrganization) ToEntity() entities.Organization {
	return entities.Organization{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address.ToEntity(),
	}
}
