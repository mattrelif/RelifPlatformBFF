package clients

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3(config aws.Config) *s3.Client {
	return s3.NewFromConfig(config)
}
