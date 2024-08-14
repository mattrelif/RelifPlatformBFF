package models

import "relif/platform-bff/entities"

type FindLocation struct {
	ID   string `bson:"id,omitempty"`
	Name string `bson:"name,omitempty"`
	Type string `bson:"type,omitempty"`
}

func (location *FindLocation) ToEntity() entities.Location {
	return entities.Location{
		ID:   location.ID,
		Name: location.Name,
		Type: location.Type,
	}
}

type Location struct {
	ID   string `bson:"id,omitempty"`
	Type string `bson:"type,omitempty"`
}

func (location *Location) ToEntity() entities.Location {
	return entities.Location{
		ID:   location.ID,
		Type: location.Type,
	}
}

func NewLocation(entity entities.Location) Location {
	return Location{
		ID:   entity.ID,
		Type: entity.Type,
	}
}
