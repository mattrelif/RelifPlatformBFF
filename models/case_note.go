package models

import (
	"relif/platform-bff/entities"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseNote struct {
	ID          string    `bson:"_id,omitempty"`
	CaseID      string    `bson:"case_id,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Content     string    `bson:"content,omitempty"`
	Tags        []string  `bson:"tags,omitempty"`
	NoteType    string    `bson:"note_type,omitempty"`
	IsImportant bool      `bson:"is_important,omitempty"`
	CreatedByID string    `bson:"created_by_id,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	UpdatedAt   time.Time `bson:"updated_at,omitempty"`
}

func (cn *CaseNote) ToEntity() entities.CaseNote {
	return entities.CaseNote{
		ID:          cn.ID,
		CaseID:      cn.CaseID,
		Title:       cn.Title,
		Content:     cn.Content,
		Tags:        cn.Tags,
		NoteType:    cn.NoteType,
		IsImportant: cn.IsImportant,
		CreatedByID: cn.CreatedByID,
		CreatedAt:   cn.CreatedAt,
		UpdatedAt:   cn.UpdatedAt,
	}
}

func NewCaseNote(entity entities.CaseNote) CaseNote {
	return CaseNote{
		ID:          primitive.NewObjectID().Hex(),
		CaseID:      entity.CaseID,
		Title:       entity.Title,
		Content:     entity.Content,
		Tags:        entity.Tags,
		NoteType:    entity.NoteType,
		IsImportant: entity.IsImportant,
		CreatedByID: entity.CreatedByID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewCaseNoteFromEntity(entity entities.CaseNote) *CaseNote {
	return &CaseNote{
		ID:          primitive.NewObjectID().Hex(),
		CaseID:      entity.CaseID,
		Title:       entity.Title,
		Content:     entity.Content,
		Tags:        entity.Tags,
		NoteType:    entity.NoteType,
		IsImportant: entity.IsImportant,
		CreatedByID: entity.CreatedByID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewUpdatedCaseNote(entity entities.CaseNote) CaseNote {
	return CaseNote{
		Title:       entity.Title,
		Content:     entity.Content,
		Tags:        entity.Tags,
		NoteType:    entity.NoteType,
		IsImportant: entity.IsImportant,
		UpdatedAt:   time.Now(),
	}
}
