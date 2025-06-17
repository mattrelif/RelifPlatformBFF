package models

import (
	"relif/platform-bff/entities"
	"time"
)

type CaseDocument struct {
	ID           string    `bson:"_id,omitempty"`
	CaseID       string    `bson:"case_id,omitempty"`
	DocumentName string    `bson:"document_name,omitempty"`
	FileName     string    `bson:"file_name,omitempty"`
	DocumentType string    `bson:"document_type,omitempty"`
	FileSize     int64     `bson:"file_size,omitempty"`
	MimeType     string    `bson:"mime_type,omitempty"`
	Description  string    `bson:"description,omitempty"`
	Tags         []string  `bson:"tags,omitempty"`
	UploadedByID string    `bson:"uploaded_by_id,omitempty"`
	FilePath     string    `bson:"file_path,omitempty"`
	DownloadURL  string    `bson:"download_url,omitempty"`
	CreatedAt    time.Time `bson:"created_at,omitempty"`
}

func (cd *CaseDocument) ToEntity() entities.CaseDocument {
	return entities.CaseDocument{
		ID:           cd.ID,
		CaseID:       cd.CaseID,
		DocumentName: cd.DocumentName,
		FileName:     cd.FileName,
		DocumentType: cd.DocumentType,
		FileSize:     cd.FileSize,
		MimeType:     cd.MimeType,
		Description:  cd.Description,
		Tags:         cd.Tags,
		UploadedByID: cd.UploadedByID,
		FilePath:     cd.FilePath,
		DownloadURL:  cd.DownloadURL,
		CreatedAt:    cd.CreatedAt,
	}
}

func NewCaseDocument(entity entities.CaseDocument) CaseDocument {
	return CaseDocument{
		CaseID:       entity.CaseID,
		DocumentName: entity.DocumentName,
		FileName:     entity.FileName,
		DocumentType: entity.DocumentType,
		FileSize:     entity.FileSize,
		MimeType:     entity.MimeType,
		Description:  entity.Description,
		Tags:         entity.Tags,
		UploadedByID: entity.UploadedByID,
		FilePath:     entity.FilePath,
		DownloadURL:  entity.DownloadURL,
		CreatedAt:    time.Now(),
	}
}

func NewCaseDocumentFromEntity(entity entities.CaseDocument) *CaseDocument {
	return &CaseDocument{
		CaseID:       entity.CaseID,
		DocumentName: entity.DocumentName,
		FileName:     entity.FileName,
		DocumentType: entity.DocumentType,
		FileSize:     entity.FileSize,
		MimeType:     entity.MimeType,
		Description:  entity.Description,
		Tags:         entity.Tags,
		UploadedByID: entity.UploadedByID,
		FilePath:     entity.FilePath,
		DownloadURL:  entity.DownloadURL,
		CreatedAt:    time.Now(),
	}
}
