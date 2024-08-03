package entities

import "time"

type Beneficiary struct {
	ID                    string
	FullName              string
	CurrentOrganizationID string
	Email                 string
	Document              Document
	Birthdate             string
	Phones                []string
	CivilStatus           string
	SpokenLanguages       []string
	Education             string
	Gender                string
	Address               Address
	Status                string
	CurrentHousingID      string
	CurrentRoomID         string
	MedicalInformation    MedicalInformation
	EmergencyContacts     []EmergencyContact
	CreatedAt             time.Time
	UpdatedAt             time.Time
	Notes                 []string
}
