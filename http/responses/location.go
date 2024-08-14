package responses

import "relif/platform-bff/entities"

type Location struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func NewLocation(entity entities.Location) Location {
	return Location{
		ID:   entity.ID,
		Name: entity.Name,
		Type: entity.Type,
	}
}
