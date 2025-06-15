package entities

import "time"

type CaseNote struct {
	ID          string
	CaseID      string
	Case        Case // Populated when needed
	Title       string
	Content     string
	Tags        []string
	NoteType    string // "CALL", "MEETING", "UPDATE", "APPOINTMENT", "OTHER"
	IsImportant bool
	CreatedByID string
	CreatedBy   User // Populated when needed
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
