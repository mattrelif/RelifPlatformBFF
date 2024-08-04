package entities

import "time"

type VoluntaryPerson struct {
	ID                 string
	OrganizationID     string
	FullName           string
	Email              string
	Documents          []Document
	Birthdate          string
	Phones             []string
	Address            Address
	Status             string
	Segments           []string
	MedicalInformation MedicalInformation
	EmergencyContacts  []EmergencyContact
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Notes              string
}
