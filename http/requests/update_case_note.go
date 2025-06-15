package requests

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UpdateCaseNoteRequest struct {
	Title       *string  `json:"title,omitempty"`
	Content     *string  `json:"content,omitempty"`
	NoteType    *string  `json:"note_type,omitempty"`
	IsImportant *bool    `json:"is_important,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (r UpdateCaseNoteRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Length(1, 255)),
		validation.Field(&r.Content, validation.Length(1, 10000)),
		validation.Field(&r.NoteType, validation.In("CALL", "MEETING", "UPDATE", "APPOINTMENT", "OTHER")),
		validation.Field(&r.Tags, validation.Each(validation.Length(1, 50))),
	)
}
