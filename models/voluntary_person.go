package models

import (
	"relif/bff/entities"
	"time"
)

type VoluntaryPerson struct {
	ID                 string             `bson:"_id"`
	OrganizationID     string             `bson:"organization_id"`
	FullName           string             `json:"full_name"`
	Email              string             `json:"email"`
	Document           Document           `json:"document"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	Address            Address            `json:"address"`
	Status             string             `json:"status"`
	Segments           []string           `bson:"segments"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	Notes              []string           `json:"notes"`
}

func (voluntary *VoluntaryPerson) ToEntity() entities.VoluntaryPerson {
	emergencyContacts := make([]entities.EmergencyContact, 0)

	for _, c := range voluntary.EmergencyContacts {
		contact := entities.EmergencyContact{
			Relationship: c.Relationship,
			FullName:     c.FullName,
			Emails:       c.Emails,
			Phones:       c.Phones,
		}

		emergencyContacts = append(emergencyContacts, contact)
	}

	return entities.VoluntaryPerson{
		ID:                 voluntary.ID,
		OrganizationID:     voluntary.OrganizationID,
		FullName:           voluntary.FullName,
		Email:              voluntary.Email,
		Document:           voluntary.Document.ToEntity(),
		Birthdate:          voluntary.Birthdate,
		Phones:             voluntary.Phones,
		Address:            voluntary.Address.ToEntity(),
		Status:             voluntary.Status,
		Segments:           voluntary.Segments,
		MedicalInformation: voluntary.MedicalInformation.ToEntity(),
		EmergencyContacts:  emergencyContacts,
		CreatedAt:          voluntary.CreatedAt,
		UpdatedAt:          voluntary.UpdatedAt,
		Notes:              voluntary.Notes,
	}
}
