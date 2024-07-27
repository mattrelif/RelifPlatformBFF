package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type Password interface {
	PasswordChangeRequest(email string) error
	UpdatePassword(id, password string) error
}

type passwordImpl struct {
	emailService   Email
	userService    Users
	repository     repositories.PasswordChangeRequests
	passwordHashFn utils.PasswordHashFn
}

func NewPassword(
	emailService Email,
	usersService Users,
	repository repositories.PasswordChangeRequests,
	passwordHashFn utils.PasswordHashFn,
) Password {
	return &passwordImpl{
		emailService:   emailService,
		userService:    usersService,
		repository:     repository,
		passwordHashFn: passwordHashFn,
	}
}

func (service *passwordImpl) PasswordChangeRequest(email string) error {
	user, err := service.userService.FindOneByEmail(email)

	if err != nil {
		return err
	}

	request := entities.PasswordChangeRequest{
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 3),
	}

	id, err := service.repository.Create(request)

	if err != nil {
		return err
	}

	if err = service.emailService.SendPasswordResetEmail(id, user); err != nil {
		return err
	}

	return nil
}

func (service *passwordImpl) UpdatePassword(id, password string) error {
	request, err := service.repository.FindOneAndDeleteById(id)

	if err != nil {
		return err
	}

	hashed, err := service.passwordHashFn(password)

	if err != nil {
		return err
	}

	data := entities.User{
		Password: hashed,
	}
	user, err := service.userService.FindOneAndUpdateById(request.UserID, data)

	if err != nil {
		return err
	}

	if err = service.emailService.SendPasswordChangedEmail(user); err != nil {
		return err
	}

	return nil
}
