package entities

import "time"

type CaseDocument struct {
	ID           string
	CaseID       string
	Case         Case // Populated when needed
	DocumentName string
	FileName     string
	DocumentType string // "FORM", "REPORT", "EVIDENCE", "CORRESPONDENCE", "IDENTIFICATION", "LEGAL", "MEDICAL", "OTHER"
	FileSize     int64
	MimeType     string
	Description  string
	Tags         []string
	UploadedByID string
	UploadedBy   User // Populated when needed
	DownloadURL  string
	CreatedAt    time.Time
}
