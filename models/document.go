package models

import "relif/bff/entities"

type Document struct {
	Type  string `bson:"type"`
	Value string `bson:"value"`
}

func (document *Document) ToEntity() entities.Document {
	return entities.Document{
		Type:  document.Type,
		Value: document.Value,
	}
}
