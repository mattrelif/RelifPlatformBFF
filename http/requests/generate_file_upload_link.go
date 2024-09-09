package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type GenerateFileUploadLink struct {
	FileType string `json:"file_type"`
}

func (req *GenerateFileUploadLink) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.FileType, validation.Required),
	)
}

func (req *GenerateFileUploadLink) ToEntity() entities.File {
	return entities.File{
		Type: req.FileType,
	}
}
