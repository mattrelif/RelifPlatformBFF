package services

import (
	"relif/bff/entities"
	"relif/bff/repositories"
	"relif/bff/utils"
	"time"
)

type Password interface {
	PasswordChangeRequest(email string) error
	UpdatePassword(code, password string) error
}

type passwordImpl struct {
	emailService   Email
	userService    Users
	repository     repositories.PasswordChangeRequests
	passwordHashFn utils.PasswordHashFn
	uuidGenerator  utils.UuidGenerator
}

func NewPassword(
	emailService Email,
	usersService Users,
	repository repositories.PasswordChangeRequests,
	passwordHashFn utils.PasswordHashFn,
	uuidGenerator utils.UuidGenerator,
) Password {
	return &passwordImpl{
		emailService:   emailService,
		userService:    usersService,
		repository:     repository,
		passwordHashFn: passwordHashFn,
		uuidGenerator:  uuidGenerator,
	}
}

func (service *passwordImpl) PasswordChangeRequest(email string) error {
	user, err := service.userService.FindOneByEmail(email)

	if err != nil {
		return err
	}

	request := entities.PasswordChangeRequest{
		UserID:    user.ID,
		Code:      service.uuidGenerator(),
		ExpiresAt: time.Now().Add(time.Hour * 3),
	}

	if err = service.repository.Create(request); err != nil {
		return err
	}

	if err = service.emailService.SendPasswordResetEmail(request.Code, user); err != nil {
		return err
	}

	return nil
}

func (service *passwordImpl) UpdatePassword(code, password string) error {
	request, err := service.repository.FindOneAndDeleteByCode(code)

	if err != nil {
		return err
	}

	hashed, err := service.passwordHashFn(password)

	if err != nil {
		return err
	}

	user, err := service.userService.FindOneById(request.UserID)

	if err != nil {
		return err
	}

	data := entities.User{
		Password: hashed,
	}
	if err = service.userService.UpdateOneById(user.ID, data); err != nil {
		return err
	}

	if err = service.emailService.SendPasswordChangedEmail(user); err != nil {
		return err
	}

	return nil
}
