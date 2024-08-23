package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type UserPreferences struct {
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

func (req *UserPreferences) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Language, validation.Required),
		validation.Field(&req.Timezone, validation.Required),
	)
}

func (req *UserPreferences) ToEntity() entities.UserPreferences {
	return entities.UserPreferences{
		Language: req.Language,
		Timezone: req.Timezone,
	}
}
