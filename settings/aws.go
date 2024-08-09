package settings

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"relif/platform-bff/utils"
)

func NewAWSConfig(region string) (aws.Config, error) {
	if region == "" {
		return aws.Config{}, utils.ErrMissingAWSRegion
	}

	return config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
}
