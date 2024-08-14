package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type AllocateProductType struct {
	To       Location `json:"to"`
	Quantity int      `json:"quantity"`
}

func (req *AllocateProductType) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.To, validation.By(func(value interface{}) error {
			if location, ok := value.(Location); ok {
				return location.Validate()
			}
			return nil
		})),
		validation.Field(&req.Quantity, validation.Required),
	)
}

func (req *AllocateProductType) ToEntity() entities.ProductTypeAllocation {
	return entities.ProductTypeAllocation{
		To:       req.To.ToEntity(),
		Quantity: req.Quantity,
	}
}
