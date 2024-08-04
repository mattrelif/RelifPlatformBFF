package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type Beneficiary struct {
	ID                    string             `bson:"_id,omitempty"`
	CurrentOrganizationID string             `bson:"current_organization_id"`
	FullName              string             `bson:"full_name,omitempty"`
	Email                 string             `bson:"email,omitempty"`
	Documents             []Document         `bson:"documents,omitempty"`
	Birthdate             string             `bson:"birthdate,omitempty"`
	Phones                []string           `bson:"phones,omitempty"`
	CivilStatus           string             `bson:"civil_status,omitempty"`
	SpokenLanguages       []string           `bson:"spoken_languages,omitempty"`
	Education             string             `bson:"education,omitempty"`
	Gender                string             `bson:"gender,omitempty"`
	Occupation            string             `bson:"occupation,omitempty"`
	Address               Address            `bson:"address,omitempty"`
	Status                string             `bson:"status,omitempty"`
	CurrentHousingID      string             `bson:"current_housing_id,omitempty"`
	CurrentRoomID         string             `bson:"current_room_id,omitempty"`
	MedicalInformation    MedicalInformation `bson:"medical_information,omitempty"`
	EmergencyContacts     []EmergencyContact `bson:"emergency_contacts,omitempty"`
	CreatedAt             time.Time          `bson:"created_at,omitempty"`
	UpdatedAt             time.Time          `bson:"updated_at,omitempty"`
	Notes                 string             `bson:"notes,omitempty"`
}

func (beneficiary *Beneficiary) ToEntity() entities.Beneficiary {
	emergencyContacts := make([]entities.EmergencyContact, 0)

	for _, contact := range beneficiary.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, contact.ToEntity())
	}

	documents := make([]entities.Document, 0)

	for _, document := range beneficiary.Documents {
		documents = append(documents, document.ToEntity())
	}

	return entities.Beneficiary{
		ID:                    beneficiary.ID,
		CurrentOrganizationID: beneficiary.CurrentOrganizationID,
		FullName:              beneficiary.FullName,
		Email:                 beneficiary.Email,
		Documents:             documents,
		Birthdate:             beneficiary.Birthdate,
		Phones:                beneficiary.Phones,
		CivilStatus:           beneficiary.CivilStatus,
		SpokenLanguages:       beneficiary.SpokenLanguages,
		Education:             beneficiary.Education,
		Gender:                beneficiary.Gender,
		Occupation:            beneficiary.Occupation,
		Address:               beneficiary.Address.ToEntity(),
		Status:                beneficiary.Status,
		CurrentHousingID:      beneficiary.CurrentHousingID,
		CurrentRoomID:         beneficiary.CurrentRoomID,
		MedicalInformation:    beneficiary.MedicalInformation.ToEntity(),
		EmergencyContacts:     emergencyContacts,
		CreatedAt:             beneficiary.CreatedAt,
		UpdatedAt:             beneficiary.UpdatedAt,
		Notes:                 beneficiary.Notes,
	}
}

func NewBeneficiary(entity entities.Beneficiary) Beneficiary {
	emergencyContacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return Beneficiary{
		ID:                    primitive.NewObjectID().Hex(),
		CurrentOrganizationID: entity.CurrentOrganizationID,
		FullName:              entity.FullName,
		Email:                 entity.Email,
		Documents:             documents,
		Birthdate:             entity.Birthdate,
		Phones:                entity.Phones,
		CivilStatus:           entity.CivilStatus,
		SpokenLanguages:       entity.SpokenLanguages,
		Education:             entity.Education,
		Gender:                entity.Gender,
		Occupation:            entity.Occupation,
		Address:               NewAddress(entity.Address),
		Status:                utils.ActiveStatus,
		MedicalInformation:    NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:     emergencyContacts,
		CreatedAt:             time.Now(),
		Notes:                 entity.Notes,
	}
}

func NewUpdatedBeneficiary(entity entities.Beneficiary) Beneficiary {
	emergencyContacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return Beneficiary{
		CurrentOrganizationID: entity.CurrentOrganizationID,
		FullName:              entity.FullName,
		Email:                 entity.Email,
		Documents:             documents,
		Birthdate:             entity.Birthdate,
		Phones:                entity.Phones,
		CivilStatus:           entity.CivilStatus,
		SpokenLanguages:       entity.SpokenLanguages,
		Education:             entity.Education,
		Gender:                entity.Gender,
		Occupation:            entity.Occupation,
		Address:               NewAddress(entity.Address),
		Status:                utils.ActiveStatus,
		MedicalInformation:    NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:     emergencyContacts,
		UpdatedAt:             time.Now(),
		Notes:                 entity.Notes,
		CurrentHousingID:      entity.CurrentHousingID,
		CurrentRoomID:         entity.CurrentRoomID,
	}
}
