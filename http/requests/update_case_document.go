package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UpdateCaseDocumentRequest struct {
	DocumentName string   `json:"document_name"`
	DocumentType string   `json:"document_type"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
}

func (r UpdateCaseDocumentRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DocumentName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.DocumentType, validation.Required, validation.In(
			"FORM", "REPORT", "EVIDENCE", "CORRESPONDENCE",
			"IDENTIFICATION", "LEGAL", "MEDICAL", "OTHER",
		)),
		validation.Field(&r.Description, validation.Length(0, 1000)),
		validation.Field(&r.Tags, validation.Each(validation.Length(1, 50))),
	)
}
