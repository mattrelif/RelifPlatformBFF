package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type OrganizationSignUp struct {
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	Email          string          `json:"email"`
	Password       string          `json:"password"`
	Phones         []string        `json:"phones"`
	Role           string          `json:"role"`
	OrganizationId string          `json:"organization_id"`
	Preferences    UserPreferences `json:"preferences"`
}

func (req *OrganizationSignUp) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required),
		validation.Field(&req.Phones, validation.Each(validation.Required)),
		validation.Field(&req.Role, validation.Required),
		validation.Field(&req.OrganizationId, validation.Required, is.MongoID),
		validation.Field(&req.Preferences, validation.By(func(value interface{}) error {
			if preferences, ok := value.(UserPreferences); ok {
				return preferences.Validate()
			}

			return nil
		})),
	)
}

func (req *OrganizationSignUp) ToEntity() entities.User {
	return entities.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Password:       req.Password,
		Phones:         req.Phones,
		Role:           req.Role,
		OrganizationID: req.OrganizationId,
		Preferences: entities.UserPreferences{
			Language: req.Preferences.Language,
			Timezone: req.Preferences.Timezone,
		},
	}
}
