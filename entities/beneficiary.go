package entities

import "time"

type BeneficiaryMedicalInformation struct {
	Allergies                  []string
	CurrentMedications         []string
	RecurrentMedicalConditions []string
	HealthInsurancePlans       []string
	BloodType                  string
	TakenVaccines              []string
	MentalHealthHistory        []string
	Height                     int
	Weight                     int
	CigarettesUsage            bool
	AlcoholConsumption         bool
	Disabilities               []string
}

type BeneficiaryEmergencyContactInformation struct {
	Relationship string
	FullName     string
	Emails       []string
	Phones       []string
}

type BeneficiaryExtraInformation struct {
	CivilStatus     string
	SpokenLanguages []string
	Education       string
}

type Beneficiary struct {
	ID                 string
	FullName           string
	Email              string
	Document           Document
	Birthdate          string
	Phones             []string
	Address            Address
	Status             string
	MedicalInformation BeneficiaryMedicalInformation
	EmergencyContacts  []BeneficiaryEmergencyContactInformation
	ExtraInformation   BeneficiaryExtraInformation
	CreatedAt          time.Time
	UpdatedAt          time.Time
	Notes              []string
}
