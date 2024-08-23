package password_recovery

import (
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type Change interface {
	Execute(code, newPassword string) error
}

type changeImpl struct {
	usersRepository                  repositories.Users
	passwordChangeRequestsRepository repositories.PasswordChangeRequests
	emailService                     services.Email
	passwordHashFunction             utils.PasswordHashFn
}

func NewChange(
	usersRepository repositories.Users,
	passwordChangeRequestsRepository repositories.PasswordChangeRequests,
	emailService services.Email,
	passwordHashFunction utils.PasswordHashFn,
) Change {
	return &changeImpl{
		usersRepository:                  usersRepository,
		passwordChangeRequestsRepository: passwordChangeRequestsRepository,
		emailService:                     emailService,
		passwordHashFunction:             passwordHashFunction,
	}
}

func (uc *changeImpl) Execute(code, newPassword string) error {
	request, err := uc.passwordChangeRequestsRepository.FindOneByCode(code)

	if err != nil {
		return err
	}

	hashed, err := uc.passwordHashFunction(newPassword)

	if err != nil {
		return err
	}

	user, err := uc.usersRepository.FindOneByID(request.UserID)

	if err != nil {
		return err
	}

	user.Password = hashed

	if err = uc.usersRepository.UpdateOneByID(user.ID, user); err != nil {
		return err
	}

	if err = uc.passwordChangeRequestsRepository.DeleteByCode(request.Code); err != nil {
		return err
	}

	if err = uc.emailService.SendPasswordChangedEmail(user); err != nil {
		return err
	}

	return nil
}
