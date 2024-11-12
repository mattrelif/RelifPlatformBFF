package requests

import "errors"

type Verify struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

func (req *Verify) Validate() error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Code == "" {
		return errors.New("verification code is required")
	}
	return nil
}
