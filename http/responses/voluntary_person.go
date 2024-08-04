package responses

import (
	"relif/bff/entities"
	"time"
)

type VoluntaryPeople []VoluntaryPerson

type VoluntaryPerson struct {
	ID                 string             `json:"id"`
	FullName           string             `json:"full_name"`
	Email              string             `json:"email"`
	Documents          []Document         `json:"documents"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	Address            Address            `json:"address"`
	OrganizationID     string             `json:"organization_id"`
	Segments           []string           `json:"segments"`
	Gender             string             `json:"gender"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	Notes              string             `json:"notes"`
}

func NewVoluntaryPerson(entity entities.VoluntaryPerson) VoluntaryPerson {
	contacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		contacts = append(contacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return VoluntaryPerson{
		ID:                 entity.ID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Documents:          documents,
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		OrganizationID:     entity.OrganizationID,
		Segments:           entity.Segments,
		Gender:             entity.Gender,
		EmergencyContacts:  contacts,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		Notes:              entity.Notes,
	}
}

func NewVoluntaryPeople(entityList []entities.VoluntaryPerson) VoluntaryPeople {
	res := make(VoluntaryPeople, 0)

	for _, entity := range entityList {
		res = append(res, NewVoluntaryPerson(entity))
	}

	return res
}
