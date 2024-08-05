package entities

import "time"

type Beneficiary struct {
	ID                    string
	FullName              string
	Email                 string
	Documents             []Document
	Birthdate             string
	Phones                []string
	CivilStatus           string
	SpokenLanguages       []string
	Education             string
	Gender                string
	Occupation            string
	Address               Address
	Status                string
	CurrentHousingID      string
	CurrentHousing        Housing
	CurrentRoomID         string
	CurrentRoom           HousingRoom
	CurrentOrganizationID string
	CurrentOrganization   Organization
	MedicalInformation    MedicalInformation
	EmergencyContacts     []EmergencyContact
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Notes                 string
}
