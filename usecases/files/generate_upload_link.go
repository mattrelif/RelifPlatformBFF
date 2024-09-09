package files

import (
	"fmt"
	"relif/platform-bff/entities"
	"relif/platform-bff/services"
	"relif/platform-bff/utils"
)

type GenerateUploadLink interface {
	Execute(domain, fileType string) (string, error)
}

type generateUploadLinkImpl struct {
	fileUploadsService services.FileUploads
	uuidGenerator      utils.UuidGenerator
}

func NewGenerateUploadLink(
	fileUploadsService services.FileUploads,
	uuidGenerator utils.UuidGenerator,
) GenerateUploadLink {
	return &generateUploadLinkImpl{
		fileUploadsService: fileUploadsService,
		uuidGenerator:      uuidGenerator,
	}
}

func (uc *generateUploadLinkImpl) Execute(domain, fileType string) (string, error) {
	fileName := uc.uuidGenerator()

	file := entities.File{
		Key:  fmt.Sprintf("%s/%s", domain, fileName),
		Type: fileType,
	}

	url, err := uc.fileUploadsService.GenerateUploadLink(file)

	if err != nil {
		return "", err
	}

	return url, nil
}
