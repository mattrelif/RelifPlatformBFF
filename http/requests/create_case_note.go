package requests

import (
	"relif/platform-bff/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateCaseNoteRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	NoteType    string   `json:"note_type"`
	IsImportant bool     `json:"is_important"`
	Tags        []string `json:"tags"`
}

func (r CreateCaseNoteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Content, validation.Required, validation.Length(1, 10000)),
		validation.Field(&r.NoteType, validation.Required, validation.In("CALL", "MEETING", "UPDATE", "APPOINTMENT", "OTHER")),
		validation.Field(&r.Tags, validation.Each(validation.Length(1, 50))),
	)
}

func (req *CreateCaseNoteRequest) ToEntity(caseID, createdByID string) entities.CaseNote {
	return entities.CaseNote{
		CaseID:      caseID,
		Title:       req.Title,
		Content:     req.Content,
		Tags:        req.Tags,
		NoteType:    req.NoteType,
		IsImportant: req.IsImportant,
		CreatedByID: createdByID,
	}
}
