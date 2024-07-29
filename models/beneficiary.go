package models

import (
	"relif/bff/entities"
	"time"
)

type Beneficiary struct {
	ID                 string             `bson:"_id"`
	FullName           string             `bson:"full_name"`
	Email              string             `bson:"email"`
	Document           Document           `bson:"document"`
	Birthdate          string             `bson:"birthdate"`
	Phones             []string           `bson:"phones"`
	CivilStatus        string             `bson:"civil_status"`
	SpokenLanguages    []string           `bson:"spoken_languages"`
	Education          string             `bson:"education"`
	Address            Address            `bson:"address"`
	Status             string             `bson:"status"`
	CurrentHousingID   string             `bson:"current_housing_id"`
	CurrentRoomID      string             `bson:"current_room_id"`
	MedicalInformation MedicalInformation `bson:"medical_information"`
	EmergencyContacts  []EmergencyContact `bson:"emergency_contacts"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
	Notes              []string           `bson:"notes"`
}

func (beneficiary *Beneficiary) ToEntity() entities.Beneficiary {
	emergencyContacts := make([]entities.EmergencyContact, len(beneficiary.EmergencyContacts))

	for _, c := range beneficiary.EmergencyContacts {
		contact := entities.EmergencyContact{
			Relationship: c.Relationship,
			FullName:     c.FullName,
			Emails:       c.Emails,
			Phones:       c.Phones,
		}

		emergencyContacts = append(emergencyContacts, contact)
	}

	return entities.Beneficiary{
		ID:               beneficiary.ID,
		FullName:         beneficiary.FullName,
		Email:            beneficiary.Email,
		Document:         beneficiary.Document.ToEntity(),
		Birthdate:        beneficiary.Birthdate,
		Phones:           beneficiary.Phones,
		CivilStatus:      beneficiary.CivilStatus,
		SpokenLanguages:  beneficiary.SpokenLanguages,
		Education:        beneficiary.Education,
		Address:          beneficiary.Address.ToEntity(),
		Status:           beneficiary.Status,
		CurrentHousingID: beneficiary.CurrentHousingID,
		CurrentRoomID:    beneficiary.CurrentRoomID,
		MedicalInformation: entities.MedicalInformation{
			Allergies:                  beneficiary.MedicalInformation.Allergies,
			CurrentMedications:         beneficiary.MedicalInformation.CurrentMedications,
			RecurrentMedicalConditions: beneficiary.MedicalInformation.RecurrentMedicalConditions,
			HealthInsurancePlans:       beneficiary.MedicalInformation.HealthInsurancePlans,
			BloodType:                  beneficiary.MedicalInformation.BloodType,
			TakenVaccines:              beneficiary.MedicalInformation.TakenVaccines,
			MentalHealthHistory:        beneficiary.MedicalInformation.MentalHealthHistory,
			Height:                     beneficiary.MedicalInformation.Height,
			Weight:                     beneficiary.MedicalInformation.Weight,
			CigarettesUsage:            beneficiary.MedicalInformation.CigarettesUsage,
			AlcoholConsumption:         beneficiary.MedicalInformation.AlcoholConsumption,
			Disabilities:               beneficiary.MedicalInformation.Disabilities,
		},
		EmergencyContacts: emergencyContacts,
		CreatedAt:         beneficiary.CreatedAt,
		UpdatedAt:         beneficiary.UpdatedAt,
		Notes:             beneficiary.Notes,
	}
}
