package settings

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAWSConfig(region string) (aws.Config, error) {
	// Use default region if not provided
	if region == "" {
		region = "us-east-1"
		fmt.Printf("Warning: AWS_REGION not set, defaulting to %s\n", region)
	}

	// Try to load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		fmt.Printf("Warning: Failed to load full AWS config (%v), creating minimal config\n", err)
		// Create minimal config for development/testing
		return aws.Config{
			Region: region,
		}, nil
	}

	fmt.Printf("AWS configuration loaded successfully for region: %s\n", region)
	return cfg, nil
}
