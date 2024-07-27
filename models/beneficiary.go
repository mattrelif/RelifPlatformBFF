package models

import (
	"relif/bff/entities"
	"time"
)

type BeneficiaryMedicalInformation struct {
	Allergies                  []string `bson:"allergies"`
	CurrentMedications         []string `bson:"current_medications"`
	RecurrentMedicalConditions []string `bson:"recurrent_medical_conditions"`
	HealthInsurancePlans       []string `bson:"health_insurance_plans"`
	BloodType                  string   `bson:"blood_type"`
	TakenVaccines              []string `bson:"taken_vaccines"`
	MentalHealthHistory        []string `bson:"mental_health_history"`
	Height                     int      `bson:"height"`
	Weight                     int      `bson:"weight"`
	CigarettesUsage            bool     `bson:"cigarettes_usage"`
	AlcoholConsumption         bool     `bson:"alcohol_consumption"`
	Disabilities               []string `bson:"disabilities"`
}

type BeneficiaryEmergencyContactInformation struct {
	Relationship string   `bson:"relationship"`
	FullName     string   `bson:"full_name"`
	Emails       []string `bson:"emails"`
	Phones       []string `bson:"phones"`
}

type BeneficiaryExtraInformation struct {
	CivilStatus     string   `bson:"civil_status"`
	SpokenLanguages []string `bson:"spoken_languages"`
	Education       string   `bson:"education"`
}

type Beneficiary struct {
	ID                 string                                   `bson:"_id"`
	FullName           string                                   `bson:"full_name"`
	Email              string                                   `bson:"email"`
	Document           Document                                 `bson:"document"`
	Birthdate          string                                   `bson:"birthdate"`
	Phones             []string                                 `bson:"phones"`
	Address            Address                                  `bson:"address"`
	Status             string                                   `bson:"status"`
	MedicalInformation BeneficiaryMedicalInformation            `bson:"medical_information"`
	EmergencyContacts  []BeneficiaryEmergencyContactInformation `bson:"emergency_contacts"`
	ExtraInformation   BeneficiaryExtraInformation              `bson:"extra_information"`
	CreatedAt          time.Time                                `bson:"created_at"`
	UpdatedAt          time.Time                                `bson:"updated_at"`
	Notes              []string                                 `bson:"notes"`
}

func (beneficiary *Beneficiary) ToEntity() entities.Beneficiary {
	emergencyContacts := make([]entities.BeneficiaryEmergencyContactInformation, len(beneficiary.EmergencyContacts))

	for _, c := range beneficiary.EmergencyContacts {
		contact := entities.BeneficiaryEmergencyContactInformation{
			Relationship: c.Relationship,
			FullName:     c.FullName,
			Emails:       c.Emails,
			Phones:       c.Phones,
		}

		emergencyContacts = append(emergencyContacts, contact)
	}

	return entities.Beneficiary{
		ID:        beneficiary.ID,
		FullName:  beneficiary.FullName,
		Email:     beneficiary.Email,
		Document:  beneficiary.Document.ToEntity(),
		Birthdate: beneficiary.Birthdate,
		Phones:    beneficiary.Phones,
		Address:   beneficiary.Address.ToEntity(),
		Status:    beneficiary.Status,
		MedicalInformation: entities.BeneficiaryMedicalInformation{
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
		ExtraInformation: entities.BeneficiaryExtraInformation{
			CivilStatus:     beneficiary.ExtraInformation.CivilStatus,
			SpokenLanguages: beneficiary.ExtraInformation.SpokenLanguages,
			Education:       beneficiary.ExtraInformation.Education,
		},
		CreatedAt: beneficiary.CreatedAt,
		UpdatedAt: beneficiary.UpdatedAt,
		Notes:     beneficiary.Notes,
	}
}
