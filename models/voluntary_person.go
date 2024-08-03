package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type VoluntaryPerson struct {
	ID                 string             `bson:"_id,omitempty"`
	OrganizationID     string             `bson:"organization_id,omitempty"`
	FullName           string             `bson:"full_name,omitempty"`
	Email              string             `bson:"email,omitempty"`
	Document           Document           `bson:"document,omitempty"`
	Birthdate          string             `bson:"birthdate,omitempty"`
	Phones             []string           `bson:"phones,omitempty"`
	Address            Address            `bson:"address,omitempty"`
	Status             string             `bson:"status,omitempty"`
	Segments           []string           `bson:"segments,omitempty"`
	MedicalInformation MedicalInformation `bson:"medical_information,omitempty"`
	EmergencyContacts  []EmergencyContact `bson:"emergency_contacts,omitempty"`
	CreatedAt          time.Time          `bson:"created_at,omitempty"`
	UpdatedAt          time.Time          `bson:"updated_at,omitempty"`
	Notes              []string           `bson:"notes,omitempty"`
}

func (voluntary *VoluntaryPerson) ToEntity() entities.VoluntaryPerson {
	emergencyContacts := make([]entities.EmergencyContact, 0)

	for _, contact := range voluntary.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, contact.ToEntity())
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

func NewVoluntaryPerson(entity entities.VoluntaryPerson) VoluntaryPerson {
	emergencyContacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, NewEmergencyContact(contact))
	}

	return VoluntaryPerson{
		ID:                 primitive.NewObjectID().Hex(),
		OrganizationID:     entity.OrganizationID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Document:           NewDocument(entity.Document),
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		Status:             utils.ActiveStatus,
		Segments:           entity.Segments,
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

	return VoluntaryPerson{
		OrganizationID:     entity.OrganizationID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Document:           NewDocument(entity.Document),
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		Address:            NewAddress(entity.Address),
		Status:             utils.ActiveStatus,
		Segments:           entity.Segments,
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:  emergencyContacts,
		UpdatedAt:          time.Now(),
		Notes:              entity.Notes,
	}
}
