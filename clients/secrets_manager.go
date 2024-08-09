package clients

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func NewSecretsManager(config aws.Config) *secretsmanager.Client {
	return secretsmanager.NewFromConfig(config)
}
