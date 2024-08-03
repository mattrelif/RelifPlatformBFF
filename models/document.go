package models

import "relif/bff/entities"

type Document struct {
	Type  string `bson:"type,omitempty"`
	Value string `bson:"value,omitempty"`
}

func (document *Document) ToEntity() entities.Document {
	return entities.Document{
		Type:  document.Type,
		Value: document.Value,
	}
}

func NewDocument(entity entities.Document) Document {
	return Document{
		Type:  entity.Type,
		Value: entity.Value,
	}
}
