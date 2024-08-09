package models

import "relif/platform-bff/entities"

type EmergencyContact struct {
	Relationship string   `bson:"relationship,omitempty"`
	FullName     string   `bson:"full_name,omitempty"`
	Emails       []string `bson:"emails,omitempty"`
	Phones       []string `bson:"phones,omitempty"`
}

func (contact EmergencyContact) ToEntity() entities.EmergencyContact {
	return entities.EmergencyContact{
		Relationship: contact.Relationship,
		FullName:     contact.FullName,
		Emails:       contact.Emails,
		Phones:       contact.Phones,
	}
}

func NewEmergencyContact(entity entities.EmergencyContact) EmergencyContact {
	return EmergencyContact{
		Relationship: entity.Relationship,
		FullName:     entity.FullName,
		Emails:       entity.Emails,
		Phones:       entity.Phones,
	}
}
