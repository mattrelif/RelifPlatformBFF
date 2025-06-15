package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateCaseDocumentRequest struct {
	DocumentName string   `json:"document_name"`
	DocumentType string   `json:"document_type"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
	FileName     string   `json:"file_name"`
	FileSize     int64    `json:"file_size"`
	MimeType     string   `json:"mime_type"`
	FileData     []byte   `json:"-"` // File data, not JSON serialized
}

func (r CreateCaseDocumentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DocumentName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.DocumentType, validation.Required, validation.In(
			"FORM", "REPORT", "EVIDENCE", "CORRESPONDENCE",
			"IDENTIFICATION", "LEGAL", "MEDICAL", "OTHER",
		)),
		validation.Field(&r.Description, validation.Length(0, 1000)),
		validation.Field(&r.FileName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.MimeType, validation.Required),
		validation.Field(&r.FileSize, validation.Required, validation.Min(1), validation.Max(10*1024*1024)), // 10MB max
		validation.Field(&r.Tags, validation.Each(validation.Length(1, 50))),
	)
}
