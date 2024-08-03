package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UpdatePassword struct {
	NewPassword string `json:"new_password"`
}

func (req *UpdatePassword) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.NewPassword, validation.Required),
	)
}
