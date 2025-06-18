package requests

import (
	"relif/platform-bff/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type EmergencyContact struct {
	Relationship string   `json:"relationship"`
	FullName     string   `json:"full_name"`
	Emails       []string `json:"emails"`
	Phones       []string `json:"phones"`
}

func (req *EmergencyContact) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Relationship, validation.Required),
		validation.Field(&req.FullName, validation.Required),
		validation.Field(&req.Emails, validation.Each(validation.Required, is.Email)),
		validation.Field(&req.Phones, validation.Each(validation.Required)),
		// Require at least one email or phone
		validation.Field(&req, validation.By(func(value interface{}) error {
			if len(req.Emails) == 0 && len(req.Phones) == 0 {
				return validation.NewError("contact_required", "At least one email or phone number is required")
			}
			return nil
		})),
	)
}

func (req *EmergencyContact) ToEntity() entities.EmergencyContact {
	return entities.EmergencyContact{
		Relationship: req.Relationship,
		FullName:     req.FullName,
		Emails:       req.Emails,
		Phones:       req.Phones,
	}
}
