package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type UpdateUser struct {
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	Email        string          `json:"email"`
	Phones       []string        `json:"phones"`
	Role         string          `json:"role"`
	PlatformRole string          `json:"platform_role"`
	Preferences  UserPreferences `json:"preferences"`
}

func (req *UpdateUser) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Role, validation.Required),
		validation.Field(&req.Phones, validation.Required),
		validation.Field(&req.PlatformRole, validation.Required),
		validation.Field(&req.Preferences, validation.By(func(value interface{}) error {
			if preferences, ok := value.(UserPreferences); ok {
				return preferences.Validate()
			}

			return nil
		})),
	)
}

func (req *UpdateUser) ToEntity() entities.User {
	return entities.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Phones:       req.Phones,
		Role:         req.Role,
		PlatformRole: req.PlatformRole,
		Preferences: entities.UserPreferences{
			Language: req.Preferences.Language,
			Timezone: req.Preferences.Timezone,
		},
	}
}
