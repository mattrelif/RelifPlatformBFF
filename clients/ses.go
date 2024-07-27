package clients

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

func NewSESClient(awsConfig aws.Config) *ses.Client {
	return ses.NewFromConfig(awsConfig)
}
