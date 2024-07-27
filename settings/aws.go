package settings

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSConfig(region string) (aws.Config, error) {
	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
}
