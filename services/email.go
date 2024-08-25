package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"relif/platform-bff/entities"
)

type Email interface {
	SendPasswordResetEmail(request entities.PasswordChangeRequest, user entities.User) error
	SendPasswordChangedEmail(user entities.User) error
	SendPlatformInviteEmail(inviter entities.User, invite entities.JoinPlatformInvite) error
}

type emailSes struct {
	client         *ses.Client
	emailDomain    string
	frontendDomain string
}

func NewSesEmail(client *ses.Client, emailDomain, frontendDomain string) Email {
	return &emailSes{
		client:         client,
		emailDomain:    emailDomain,
		frontendDomain: frontendDomain,
	}
}

func (service *emailSes) SendPasswordResetEmail(request entities.PasswordChangeRequest, user entities.User) error {
	templateData := map[string]string{
		"first_name": user.FirstName,
		"reset_link": fmt.Sprintf(`"https://%s/recover-password?code=%s"`, service.frontendDomain, request.Code),
	}

	templateDataJson, err := json.Marshal(templateData)

	if err != nil {
		return err
	}

	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.emailDomain),
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Template:     aws.String("PasswordResetTemplate"),
		TemplateData: aws.String(string(templateDataJson)),
	}

	if _, err = service.client.SendTemplatedEmail(context.Background(), input); err != nil {
		return err
	}

	return nil
}

func (service *emailSes) SendPasswordChangedEmail(user entities.User) error {
	templateData := map[string]string{
		"first_name": user.FirstName,
	}

	templateDataJson, err := json.Marshal(templateData)

	if err != nil {
		return err
	}

	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.emailDomain),
		Destination: &types.Destination{
			ToAddresses: []string{user.Email},
		},
		Template:     aws.String("PasswordChangeConfirmationTemplate"),
		TemplateData: aws.String(string(templateDataJson)),
	}

	if _, err = service.client.SendTemplatedEmail(context.Background(), input); err != nil {
		return err
	}

	return nil
}

func (service *emailSes) SendPlatformInviteEmail(inviter entities.User, invite entities.JoinPlatformInvite) error {
	templateData := map[string]string{
		"inviter_name":      inviter.FirstName,
		"registration_link": fmt.Sprintf("https://%s/join?code=%s", service.frontendDomain, invite.Code),
	}

	templateDataJson, err := json.Marshal(templateData)

	if err != nil {
		return err
	}

	input := &ses.SendTemplatedEmailInput{
		Source: aws.String(service.emailDomain),
		Destination: &types.Destination{
			ToAddresses: []string{invite.InvitedEmail},
		},
		Template:     aws.String("PlatformInvitationTemplate"),
		TemplateData: aws.String(string(templateDataJson)),
	}

	if _, err = service.client.SendTemplatedEmail(context.Background(), input); err != nil {
		return err
	}

	return nil
}
