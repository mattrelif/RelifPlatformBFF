package authentication

import (
	"relif/platform-bff/services"
)

type Verify interface {
	Execute(username, code string) error
}

type verifyImpl struct {
	cognitoService services.Cognito
}

func NewVerify(cognitoService services.Cognito) Verify {
	return &verifyImpl{
		cognitoService: cognitoService,
	}
}

func (uc *verifyImpl) Execute(username, code string) error {
	return uc.cognitoService.ConfirmEmailOrPhoneVerification(username, code)
}
