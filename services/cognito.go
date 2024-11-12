package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"relif/platform-bff/repositories"
	"relif/platform-bff/settings"
	"relif/platform-bff/utils"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"
)

var (
	settingsInstance *settings.Settings
)

// Cognito is the interface for the Cognito service.
type Cognito interface {
	InitiateEmailOrPhoneVerification(username string, password string, name string, locale string) error
	ConfirmEmailOrPhoneVerification(username, code string) error
	ResendConfirmationCode(username string) error
	AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error)
}

// cognitoImpl is the implementation of the Cognito interface.
type cognitoImpl struct {
	client          *cognitoidentityprovider.Client
	clientID        string
	usersRepository repositories.Users
}

// NewCognito creates a new instance of the Cognito service.
func NewCognito(
	region string,
	clientID string,
	usersRepository repositories.Users,
) (Cognito, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	client := cognitoidentityprovider.NewFromConfig(cfg)
	return &cognitoImpl{
		client:          client,
		clientID:        clientID,
		usersRepository: usersRepository,
	}, nil
}

func (c *cognitoImpl) AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(settingsInstance.POOL_ID),
		Username:   aws.String(username),
	}
	return c.client.AdminGetUser(context.TODO(), input)
}

// InitiateEmailOrPhoneVerification sends a confirmation code to the user's email or phone.
func (c *cognitoImpl) InitiateEmailOrPhoneVerification(username string, password string, name string, locale string) error {
	log.Printf("Locale is: %s", locale)
	getUserInput := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(settingsInstance.POOL_ID), //Needs to set as Env variable
		Username:   aws.String(username),
	}

	// Attempt to get the user
	userResp, err := c.client.AdminGetUser(context.TODO(), getUserInput)
	if err == nil {
		// User exists, check if they need confirmation
		if userResp.UserStatus == types.UserStatusTypeConfirmed {
			log.Printf("User %s is already confirmed. No need to resend verification code.", username)
			return nil
		}

		log.Printf("User %s already exists and is not confirmed. Resending confirmation code.", username)

		// Resend confirmation code if the user exists and is not confirmed
		resendInput := &cognitoidentityprovider.ResendConfirmationCodeInput{
			ClientId: aws.String(c.clientID),
			Username: aws.String(username),
			ClientMetadata: map[string]string{
				"locale": locale,
			},
		}

		_, err = c.client.ResendConfirmationCode(context.TODO(), resendInput)
		if err != nil {
			var apiErr smithy.APIError
			if errors.As(err, &apiErr) {
				log.Printf("Cognito API error: Code=%s, Message=%s, Fault=%s", apiErr.ErrorCode(), apiErr.ErrorMessage(), apiErr.ErrorFault())
			}
			return fmt.Errorf("failed to resend confirmation code: %w", err)
		}

		log.Printf("Verification code resent for user: %s", username)
		return nil
	}

	// Prepare sign up input
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(username),
		Password: aws.String(password), // You should generate a secure temporary password
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(username), // Assuming username is the email
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(name), // Use the provided name
			},
		},
		ClientMetadata: map[string]string{
			"locale": locale,
		},
	}

	// Attempt to sign up the user
	_, err = c.client.SignUp(context.TODO(), signUpInput)
	if err != nil {
		// Check if it's a UsernameExistsException by comparing the error code
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			if apiErr.ErrorCode() == "UsernameExistsException" {
				// Handle user already exists case
				log.Printf("User %s already exists, resending confirmation code.", username)
				return nil
			} else {
				return fmt.Errorf("failed to sign up user: %s - %w", apiErr.ErrorMessage(), err)
			}
		} else {
			// If it's a different error that doesn't match the APIError interface
			return fmt.Errorf("unknown error occurred during sign up: %w", err)
		}
	}

	log.Printf("User %s signed up successfully.", username)
	return nil
}

func (c *cognitoImpl) ResendConfirmationCode(username string) error {
	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(username),
	}

	_, err := c.client.ResendConfirmationCode(context.TODO(), input)
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			log.Printf("Cognito API error: Code=%s, Message=%s", apiErr.ErrorCode(), apiErr.ErrorMessage())
		}
		return fmt.Errorf("failed to resend confirmation code: %w", err)
	}

	log.Printf("Verification code resent for user: %s", username)
	return nil
}

// ConfirmEmailOrPhoneVerification confirms the user's email or phone using the provided code.
func (c *cognitoImpl) ConfirmEmailOrPhoneVerification(username, code string) error {
	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(code),
	}
	_, err := c.client.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		// Check if the error is the NotAuthorizedException with status already confirmed
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotAuthorizedException" && strings.Contains(apiErr.ErrorMessage(), "Current status is CONFIRMED") {
			return fmt.Errorf("account already registered: %w", err)
		}
		// Fallback for other errors
		return fmt.Errorf("failed to confirm email/phone verification: %w", err)
	}
	// After successful confirmation, update user status to Active
	user, err := c.usersRepository.FindOneByEmail(username)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	user.Status = utils.ActiveStatus
	err = c.usersRepository.UpdateOneByID(user.ID, user)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	return nil
}
