package responses

import (
	"relif/platform-bff/entities"
)

type EmergencyContact struct {
	Relationship string   `json:"relationship"`
	FullName     string   `json:"full_name"`
	Emails       []string `json:"emails"`
	Phones       []string `json:"phones"`
}

func NewEmergencyContact(entity entities.EmergencyContact) EmergencyContact {
	return EmergencyContact{
		Relationship: entity.Relationship,
		FullName:     entity.FullName,
		Emails:       entity.Emails,
		Phones:       entity.Phones,
	}
}
