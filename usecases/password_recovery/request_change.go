package password_recovery

import (
	"relif/platform-bff/entities"
	"relif/platform-bff/guards"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
	"time"
)

type RequestChange interface {
	Execute(email string) error
}

type requestChangeImpl struct {
	usersRepository                  repositories.Users
	passwordChangeRequestsRepository repositories.PasswordChangeRequests
	emailService                     services.Email
	uuidGenerator                    utils.UuidGenerator
}

func NewRequestChange(
	usersRepository repositories.Users,
	passwordChangeRequestsRepository repositories.PasswordChangeRequests,
	emailService services.Email,
	uuidGenerator utils.UuidGenerator,
) RequestChange {
	return &requestChangeImpl{
		usersRepository:                  usersRepository,
		passwordChangeRequestsRepository: passwordChangeRequestsRepository,
		emailService:                     emailService,
		uuidGenerator:                    uuidGenerator,
	}
}

func (uc *requestChangeImpl) Execute(email string) error {
	user, err := uc.usersRepository.FindOneByEmail(email)

	if err != nil {
		return err
	}

	if err = guards.CanAccessPlatform(user); err != nil {
		return err
	}

	request := entities.PasswordChangeRequest{
		UserID:    user.ID,
		Code:      uc.uuidGenerator(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}

	if err = uc.passwordChangeRequestsRepository.Create(request); err != nil {
		return err
	}

	if err = uc.emailService.SendPasswordResetEmail(request, user); err != nil {
		return err
	}

	return nil
}
