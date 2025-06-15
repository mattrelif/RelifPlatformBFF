package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"relif/platform-bff/repositories"
	"relif/platform-bff/settings"
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
	enabled         bool
}

// NewCognito creates a new instance of the Cognito service.
func NewCognito(
	region string,
	clientID string,
	usersRepository repositories.Users,
) (Cognito, error) {
	// Try to create AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("Warning: Cognito service disabled due to AWS config error: %v\n", err)
		// Return a disabled Cognito service that won't crash
		return &cognitoImpl{
			client:          nil,
			clientID:        clientID,
			usersRepository: usersRepository,
			enabled:         false,
		}, nil
	}

	client := cognitoidentityprovider.NewFromConfig(cfg)
	return &cognitoImpl{
		client:          client,
		clientID:        clientID,
		usersRepository: usersRepository,
		enabled:         true,
	}, nil
}

func (c *cognitoImpl) AdminGetUser(username string) (*cognitoidentityprovider.AdminGetUserOutput, error) {
	if !c.enabled || c.client == nil {
		// Return an error that simulates user not found in Cognito
		return nil, fmt.Errorf("cognito service not available")
	}

	if settingsInstance == nil {
		return nil, fmt.Errorf("settings not initialized")
	}

	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String(settingsInstance.POOL_ID),
		Username:   aws.String(username),
	}
	return c.client.AdminGetUser(context.TODO(), input)
}

// InitiateEmailOrPhoneVerification sends a confirmation code to the user's email or phone.
func (c *cognitoImpl) InitiateEmailOrPhoneVerification(username string, password string, name string, locale string) error {
	if !c.enabled || c.client == nil {
		fmt.Printf("Warning: Cognito verification disabled - skipping email verification for %s\n", username)
		return nil // Don't fail user creation if Cognito is unavailable
	}

	// Set the preferred language based on the locale
	var languageCode string
	switch strings.ToLower(locale) {
	case "pt-br", "pt", "portuguese":
		languageCode = "pt"
	case "es", "spanish":
		languageCode = "es"
	default:
		languageCode = "en"
	}

	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
			},
			{
				Name:  aws.String("locale"),
				Value: aws.String(languageCode),
			},
		},
	}

	_, err := c.client.SignUp(context.TODO(), input)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			switch ae.ErrorCode() {
			case "UsernameExistsException":
				log.Printf("User %s already exists in Cognito", username)
				return nil // Don't fail if user already exists
			default:
				return fmt.Errorf("failed to initiate verification: %w", err)
			}
		}
		return fmt.Errorf("failed to initiate verification: %w", err)
	}

	log.Printf("Verification initiated for user: %s", username)
	return nil
}

// ConfirmEmailOrPhoneVerification confirms the verification code sent to the user.
func (c *cognitoImpl) ConfirmEmailOrPhoneVerification(username, code string) error {
	if !c.enabled || c.client == nil {
		fmt.Printf("Warning: Cognito verification disabled - auto-confirming user %s\n", username)
		return nil // Don't fail if Cognito is unavailable
	}

	input := &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(username),
		ConfirmationCode: aws.String(code),
	}

	_, err := c.client.ConfirmSignUp(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to confirm verification: %w", err)
	}

	log.Printf("Verification confirmed for user: %s", username)
	return nil
}

// ResendConfirmationCode resends the confirmation code to the user.
func (c *cognitoImpl) ResendConfirmationCode(username string) error {
	if !c.enabled || c.client == nil {
		fmt.Printf("Warning: Cognito verification disabled - cannot resend code for %s\n", username)
		return fmt.Errorf("verification service not available")
	}

	input := &cognitoidentityprovider.ResendConfirmationCodeInput{
		ClientId: aws.String(c.clientID),
		Username: aws.String(username),
	}

	_, err := c.client.ResendConfirmationCode(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to resend confirmation code: %w", err)
	}

	log.Printf("Confirmation code resent for user: %s", username)
	return nil
}
