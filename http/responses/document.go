package responses

import "relif/bff/entities"

type Document struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewDocument(entity entities.Document) Document {
	return Document{
		Type:  entity.Type,
		Value: entity.Value,
	}
}
