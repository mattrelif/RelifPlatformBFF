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
	Document           Document           `json:"document"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	Address            Address            `json:"address"`
	OrganizationID     string             `json:"organization_id"`
	Segments           []string           `json:"segments"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	Notes              []string           `json:"notes"`
}

func NewVoluntaryPerson(entity entities.VoluntaryPerson) VoluntaryPerson {
	contactsEntityList := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		contactsEntityList = append(contactsEntityList, NewEmergencyContact(contact))
	}

	return VoluntaryPerson{
		ID:                 entity.ID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Document:           NewDocument(entity.Document),
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		OrganizationID:     entity.OrganizationID,
		Segments:           entity.Segments,
		EmergencyContacts:  contactsEntityList,
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
