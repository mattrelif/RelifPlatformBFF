package responses

import (
	"relif/bff/entities"
	"time"
)

type Beneficiaries []Beneficiary

type Beneficiary struct {
	ID                    string             `json:"id"`
	FullName              string             `json:"full_name"`
	Email                 string             `json:"email"`
	Documents             []Document         `json:"documents"`
	Birthdate             string             `json:"birthdate"`
	Phones                []string           `json:"phones"`
	CivilStatus           string             `json:"civil_status"`
	SpokenLanguages       []string           `json:"spoken_languages"`
	Education             string             `json:"education"`
	Gender                string             `json:"gender"`
	Occupation            string             `json:"occupation"`
	Address               Address            `json:"address"`
	CurrentHousingID      string             `json:"current_housing_id"`
	CurrentHousing        Housing            `json:"current_housing,omitempty"`
	CurrentRoomID         string             `json:"current_room_id"`
	CurrentRoom           HousingRoom        `json:"current_room,omitempty"`
	CurrentOrganizationID string             `json:"current_organization_id"`
	CurrentOrganization   Organization       `json:"current_organization,omitempty"`
	MedicalInformation    MedicalInformation `json:"medical_information"`
	EmergencyContacts     []EmergencyContact `json:"emergency_contacts"`
	CreatedAt             time.Time          `json:"created_at"`
	UpdatedAt             time.Time          `json:"updated_at"`
	Notes                 string             `json:"notes"`
}

func NewBeneficiary(entity entities.Beneficiary) Beneficiary {
	contacts := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		contacts = append(contacts, NewEmergencyContact(contact))
	}

	documents := make([]Document, 0)

	for _, document := range entity.Documents {
		documents = append(documents, NewDocument(document))
	}

	return Beneficiary{
		ID:                    entity.ID,
		FullName:              entity.FullName,
		Email:                 entity.Email,
		Documents:             documents,
		Birthdate:             entity.Birthdate,
		Phones:                entity.Phones,
		CivilStatus:           entity.CivilStatus,
		SpokenLanguages:       entity.SpokenLanguages,
		Education:             entity.Education,
		Gender:                entity.Gender,
		CurrentHousingID:      entity.CurrentHousingID,
		CurrentHousing:        NewHousing(entity.CurrentHousing),
		CurrentRoomID:         entity.CurrentRoomID,
		CurrentRoom:           NewHousingRoom(entity.CurrentRoom),
		CurrentOrganizationID: entity.CurrentOrganizationID,
		CurrentOrganization:   NewOrganization(entity.CurrentOrganization),
		Occupation:            entity.Occupation,
		Address:               NewAddress(entity.Address),
		MedicalInformation:    NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:     contacts,
		CreatedAt:             entity.CreatedAt,
		UpdatedAt:             entity.UpdatedAt,
		Notes:                 entity.Notes,
	}
}

func NewBeneficiaries(entityList []entities.Beneficiary) Beneficiaries {
	res := make(Beneficiaries, 0)

	for _, entity := range entityList {
		res = append(res, NewBeneficiary(entity))
	}

	return res
}
