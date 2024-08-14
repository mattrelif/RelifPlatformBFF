package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/platform-bff/entities"
)

type CreateDonation struct {
	From          Location `json:"from"`
	ProductTypeID string   `json:"product_type_id"`
	Quantity      int      `json:"quantity"`
}

func (req *CreateDonation) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.From, validation.By(func(value interface{}) error {
			if location, ok := value.(Location); ok {
				return location.Validate()
			}
			return nil
		})),
		validation.Field(&req.ProductTypeID, validation.Required, is.MongoID),
		validation.Field(&req.Quantity, validation.Required),
	)
}

func (req *CreateDonation) ToEntity() entities.Donation {
	return entities.Donation{
		From:          req.From.ToEntity(),
		ProductTypeID: req.ProductTypeID,
		Quantity:      req.Quantity,
	}
}
