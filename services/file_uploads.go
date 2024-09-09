package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"relif/platform-bff/entities"
	"time"
)

type FileUploads interface {
	GenerateUploadLink(file entities.File) (string, error)
}

type s3FileUploads struct {
	client     *s3.Client
	bucketName string
}

func NewS3FileUploads(client *s3.Client, bucketName string) FileUploads {
	return &s3FileUploads{
		client:     client,
		bucketName: bucketName,
	}
}

func (service *s3FileUploads) GenerateUploadLink(file entities.File) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(service.bucketName),
		Key:         aws.String(file.Key),
		ContentType: aws.String(file.Type),
	}

	presigner := s3.NewPresignClient(service.client)

	req, err := presigner.PresignPutObject(context.Background(), input, s3.WithPresignExpires(time.Minute))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
