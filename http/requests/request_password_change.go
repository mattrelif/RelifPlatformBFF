package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RequestPasswordChange struct {
	Email string `json:"email"`
}

func (req *RequestPasswordChange) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
	)
}
