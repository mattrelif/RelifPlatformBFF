package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type ReallocateProductType struct {
	From     Location `json:"from"`
	To       Location `json:"to"`
	Quantity int      `json:"quantity"`
}

func (req *ReallocateProductType) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.From, validation.By(func(value interface{}) error {
			if location, ok := value.(Location); ok {
				return location.Validate()
			}
			return nil
		})),
		validation.Field(&req.To, validation.By(func(value interface{}) error {
			if location, ok := value.(Location); ok {
				return location.Validate()
			}
			return nil
		})),
		validation.Field(&req.Quantity, validation.Required),
	)
}

func (req *ReallocateProductType) ToEntity() entities.ProductTypeAllocation {
	return entities.ProductTypeAllocation{
		From:     req.From.ToEntity(),
		To:       req.To.ToEntity(),
		Quantity: req.Quantity,
	}
}
