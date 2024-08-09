package settings

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kelseyhightower/envconfig"
	"relif/platform-bff/utils"
	"strings"
	"time"
)

type Settings struct {
	Cors struct {
		AllowedOrigin string `default:"http://localhost:3000" json:"CORS_ALLOWED_ORIGIN"`
	}

	Email struct {
		Domain string `required:"true" json:"EMAIL_DOMAIN"`
	}

	Router struct {
		Context string `default:"/api/v1" json:"ROUTER_CONTEXT"`
	}

	Server struct {
		Port         int           `default:"8080" json:"SERVER_PORT"`
		WriteTimeout time.Duration `default:"10s" json:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `default:"10s" json:"SERVER_READ_TIMEOUT"`
	}

	Mongo struct {
		URI               string        `default:"mongodb://127.0.0.1:27017" json:"MONGO_URI"`
		Database          string        `default:"test" json:"MONGO_DATABASE"`
		ConnectionTimeout time.Duration `default:"10s" json:"MONGO_CONNECTION_TIMEOUT"`
	}
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

		if err = json.Unmarshal([]byte(*result.SecretString), &settings); err != nil {
			return nil, err
		}
	}

	return &settings, nil
}
