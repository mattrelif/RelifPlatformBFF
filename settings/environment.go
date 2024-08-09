package settings

import "os"

var (
	Environment = os.Getenv("ENVIRONMENT")
	AWSRegion   = os.Getenv("AWS_REGION")
	SecretName  = os.Getenv("SECRET_NAME")
)
