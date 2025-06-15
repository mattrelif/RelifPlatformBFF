package settings

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kelseyhightower/envconfig"
)

type Settings struct {
	FrontendDomain string `default:"localhost:3000" split_words:"true" json:"FRONTEND_DOMAIN"`

	EmailDomain string `default:"localhost" split_words:"true" json:"EMAIL_DOMAIN"`

	RouterContext string `default:"/api/v1" split_words:"true" json:"ROUTER_CONTEXT"`

	ServerPort string `default:"8080" split_words:"true" json:"SERVER_PORT"`

	TokenSecret string `default:"development-secret-change-in-production" split_words:"true" json:"TOKEN_SECRET"`

	MongoURI      string `default:"mongodb://127.0.0.1:27017" split_words:"true" json:"MONGO_URI"`
	MongoDatabase string `default:"relif_dev" split_words:"true" json:"MONGO_DATABASE"`

	S3BucketName string `default:"relif-dev-bucket" split_words:"true" json:"S3_BUCKET_NAME"`

	AWS_REGION string `default:"us-east-1" split_words:"true" json:"AWS_REGION"`

	COGNITO_CLIENT_ID string `default:"dev-client-id" split_words:"true" json:"COGNITO_CLIENT_ID"`

	POOL_ID string `default:"dev-pool-id" split_words:"true" json:"POOL_ID"`
}

func NewSettings(secretsManagerClient *secretsmanager.Client) (*Settings, error) {
	var settings Settings

	// Default to development if ENVIRONMENT is not set
	environment := Environment
	if environment == "" {
		environment = "development"
		fmt.Println("Warning: ENVIRONMENT not set, defaulting to 'development'")
	}

	switch strings.ToLower(environment) {
	case "development", "dev", "local":
		// Use environment variables with defaults for development
		if err := envconfig.Process("", &settings); err != nil {
			return nil, fmt.Errorf("failed to process environment variables: %w", err)
		}
		fmt.Printf("Loaded development configuration (MongoDB: %s, Port: %s)\n",
			settings.MongoDatabase, settings.ServerPort)

	case "production", "prod":
		// Try to load from AWS Secrets Manager first
		if SecretName != "" && secretsManagerClient != nil {
			fmt.Printf("Loading production configuration from AWS Secrets Manager: %s\n", SecretName)

			if err := loadFromSecretsManager(secretsManagerClient, &settings); err != nil {
				fmt.Printf("Warning: Failed to load from Secrets Manager (%v), falling back to environment variables\n", err)
				// Fallback to environment variables
				if err := envconfig.Process("", &settings); err != nil {
					return nil, fmt.Errorf("failed to load configuration from both Secrets Manager and environment: %w", err)
				}
			}
		} else {
			fmt.Println("SECRET_NAME not provided or Secrets Manager client unavailable, using environment variables")
			// Use environment variables for production when Secrets Manager is not available
			if err := envconfig.Process("", &settings); err != nil {
				return nil, fmt.Errorf("failed to process environment variables for production: %w", err)
			}
		}

	default:
		return nil, fmt.Errorf("unsupported ENVIRONMENT value: %s (supported: development, production)", environment)
	}

	// Validate critical settings
	if err := validateSettings(&settings); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return &settings, nil
}

func loadFromSecretsManager(client *secretsmanager.Client, settings *Settings) error {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(SecretName),
	}

	result, err := client.GetSecretValue(context.Background(), input)
	if err != nil {
		return fmt.Errorf("failed to get secret from AWS Secrets Manager: %w", err)
	}

	secretString := *result.SecretString
	if err = json.Unmarshal([]byte(secretString), settings); err != nil {
		return fmt.Errorf("failed to parse secret JSON: %w", err)
	}

	return nil
}

func validateSettings(settings *Settings) error {
	if settings.TokenSecret == "" || settings.TokenSecret == "development-secret-change-in-production" {
		if Environment == "production" {
			return fmt.Errorf("TOKEN_SECRET must be set to a secure value in production")
		}
		fmt.Println("Warning: Using default TOKEN_SECRET - change this in production!")
	}

	if settings.MongoURI == "" {
		return fmt.Errorf("MONGO_URI is required")
	}

	if settings.ServerPort == "" {
		settings.ServerPort = "8080"
	}

	return nil
}
