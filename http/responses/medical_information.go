package responses

import (
	"relif/bff/entities"
)

type MedicalInformation struct {
	Allergies                  []string `json:"allergies"`
	CurrentMedications         []string `json:"current_medications"`
	RecurrentMedicalConditions []string `json:"recurrent_medical_conditions"`
	HealthInsurancePlans       []string `json:"health_insurance_plans"`
	BloodType                  string   `json:"blood_type"`
	TakenVaccines              []string `json:"taken_vaccines"`
	MentalHealthHistory        []string `json:"mental_health_history"`
	Height                     int      `json:"height"`
	Weight                     int      `json:"weight"`
	CigarettesUsage            bool     `json:"cigarettes_usage"`
	AlcoholConsumption         bool     `json:"alcohol_consumption"`
	Disabilities               []string `json:"disabilities"`
}

func NewMedicalInformation(entity entities.MedicalInformation) MedicalInformation {
	return MedicalInformation{
		Allergies:                  entity.Allergies,
		CurrentMedications:         entity.CurrentMedications,
		RecurrentMedicalConditions: entity.RecurrentMedicalConditions,
		HealthInsurancePlans:       entity.HealthInsurancePlans,
		BloodType:                  entity.BloodType,
		TakenVaccines:              entity.TakenVaccines,
		MentalHealthHistory:        entity.MentalHealthHistory,
		Height:                     entity.Height,
		Weight:                     entity.Weight,
		CigarettesUsage:            entity.CigarettesUsage,
		AlcoholConsumption:         entity.AlcoholConsumption,
		Disabilities:               entity.Disabilities,
	}
}
