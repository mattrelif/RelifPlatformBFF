package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"time"
)

type VoluntaryPerson struct {
	ID                 string             `bson:"_id,omitempty"`
	OrganizationID     string             `bson:"organization_id,omitempty"`
	FullName           string             `bson:"full_name,omitempty"`
	Email              string             `bson:"email,omitempty"`
	Documents          []Document         `bson:"documents,omitempty"`
	Birthdate          string             `bson:"birthdate,omitempty"`
	Phones             []string           `bson:"phones,omitempty"`
	Address            Address            `bson:"address,omitempty"`
	Status             string             `bson:"status,omitempty"`
	Segments           []string           `bson:"segments,omitempty"`
	Gender             string             `bson:"gender,omitempty"`
	MedicalInformation MedicalInformation `bson:"medical_information,omitempty"`
	EmergencyContacts  []EmergencyContact `bson:"emergency_contacts,omitempty"`
	CreatedAt          time.Time          `bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `bson:"updated_at,omitempty"`
	Notes              string             `bson:"notes,omitempty"`
}

func (voluntary *VoluntaryPerson) ToEntity() entities.VoluntaryPerson {
	emergencyContacts := make([]entities.EmergencyContact, 0)

	for _, contact := range voluntary.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, contact.ToEntity())
	}

	documents := make([]entities.Document, 0)

	for _, document := range voluntary.Documents {
		documents = append(documents, document.ToEntity())
	}

	return entities.VoluntaryPerson{
		ID:                 voluntary.ID,
		OrganizationID:     voluntary.OrganizationID,
		FullName:           voluntary.FullName,
		Email:              voluntary.Email,
		Documents:          documents,
		Birthdate:          voluntary.Birthdate,
		Phones:             voluntary.Phones,
		Address:            voluntary.Address.ToEntity(),
		Status:             voluntary.Status,
		Segments:           voluntary.Segments,
		Gender:             voluntary.Gender,
		MedicalInformation: voluntary.MedicalInformation.ToEntity(),
		EmergencyContacts:  emergencyContacts,
		CreatedAt:          voluntary.CreatedAt,
		UpdatedAt:          voluntary.UpdatedAt,
		Notes:              voluntary.Notes,
	}
}

func NewVoluntaryPerson(entity entities.VoluntaryPerson) VoluntaryPerson {
	emergencyContacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return VoluntaryPerson{
		ID:                 primitive.NewObjectID().Hex(),
		OrganizationID:     entity.OrganizationID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Documents:          documents,
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		Status:             utils.ActiveStatus,
		Segments:           entity.Segments,
		Gender:             entity.Gender,
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:  emergencyContacts,
		CreatedAt:          time.Now(),
		Notes:              entity.Notes,
	}
}

func NewUpdatedVoluntaryPerson(entity entities.VoluntaryPerson) VoluntaryPerson {
	emergencyContacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return VoluntaryPerson{
		OrganizationID:     entity.OrganizationID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Documents:          documents,
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		Status:             utils.ActiveStatus,
		Segments:           entity.Segments,
		Gender:             entity.Gender,
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:  emergencyContacts,
		UpdatedAt:          time.Now(),
		Notes:              entity.Notes,
	}
}
