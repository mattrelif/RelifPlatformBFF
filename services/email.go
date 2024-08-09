package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"relif/bff/entities"
)

type Email interface {
	SendPasswordResetEmail(requestID string, user entities.User) error
	SendPasswordChangedEmail(user entities.User) error
	SendPlatformInviteEmail(inviter entities.User, invitedEmail, code string) error
}

type emailSes struct {
	client *ses.Client
	domain string
}

func NewSesEmail(client *ses.Client, domain string) Email {
	return &emailSes{
		client: client,
		domain: domain,
	}
}

func (service *emailSes) SendPasswordResetEmail(requestID string, user entities.User) error {
	var greetingLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	var subjectLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	var textLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.domain),
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Template: aws.String("PasswordResetEmail"),
		TemplateData: aws.String(
			fmt.Sprintf(`{"greeting": "%s", "subject": "%s", "text": "%s", "token": "%s", "first_name": "%s"}`,
				greetingLanguageMap[user.Preferences.Language],
				subjectLanguageMap[user.Preferences.Language],
				textLanguageMap[user.Preferences.Language],
				requestID,
				user.FirstName,
			)),
	}

	_, err := service.client.SendTemplatedEmail(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

func (service *emailSes) SendPasswordChangedEmail(user entities.User) error {
	var greetingLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	var subjectLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	var textLanguageMap = map[string]string{
		"pt": "",
		"en": "",
		"es": "",
	}

	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.domain),
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Template: aws.String("PasswordChangedEmail"),
		TemplateData: aws.String(
			fmt.Sprintf(`{"greeting": "%s", "subject": "%s", "text": "%s", "first_name": "%s"}`,
				greetingLanguageMap[user.Preferences.Language],
				subjectLanguageMap[user.Preferences.Language],
				textLanguageMap[user.Preferences.Language],
				user.FirstName,
			)),
	}

	_, err := service.client.SendTemplatedEmail(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}

func (service *emailSes) SendPlatformInviteEmail(inviter entities.User, invitedEmail, code string) error {
	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.domain),
		Destination: &types.Destination{
			ToAddresses: []string{invitedEmail},
		},
		Template: aws.String("PlatformInviteEmail"),
		TemplateData: aws.String(
			fmt.Sprintf(`{"invite_link": "%s"}`,
				fmt.Sprintf(""),
			)),
	}

	_, err := service.client.SendTemplatedEmail(context.Background(), input)
	if err != nil {
		return err
	}

	return nil
}
