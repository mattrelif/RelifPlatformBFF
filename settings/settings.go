package settings

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kelseyhightower/envconfig"
	"relif/platform-bff/utils"
	"strings"
)

type Settings struct {
	CorsAllowedDomain string `default:"localhost:3000" split_words:"true" json:"CORS_ALLOWED_DOMAIN"`

	EmailDomain string `required:"true" split_words:"true" json:"EMAIL_DOMAIN"`

	RouterContext string `default:"/api/v1" split_words:"true" json:"ROUTER_CONTEXT"`

	ServerPort string `default:"8080" split_words:"true" json:"SERVER_PORT"`

	MongoURI      string `default:"mongodb://127.0.0.1:27017" split_words:"true" json:"MONGO_URI"`
	MongoDatabase string `default:"test" split_words:"true" json:"MONGO_DATABASE"`
}

func NewSettings(secretsManagerClient *secretsmanager.Client) (*Settings, error) {
	var settings Settings

	if Environment == "" {
		return nil, utils.ErrMissingEnvironmentEnvVariable
	}

	switch strings.ToLower(Environment) {
	case "development":
		if err := envconfig.Process("", &settings); err != nil {
			return nil, err
		}
	case "production":
		if SecretName == "" {
			return nil, utils.ErrMissingSecretNameEnvVariable
		}

		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(SecretName),
		}

		result, err := secretsManagerClient.GetSecretValue(context.Background(), input)

		if err != nil {
			return nil, err
		}

		secretString := *result.SecretString

		if err = json.Unmarshal([]byte(secretString), &settings); err != nil {
			return nil, err
		}
	}

	return &settings, nil
}
