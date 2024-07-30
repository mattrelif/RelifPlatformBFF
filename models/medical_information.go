package models

import "relif/bff/entities"

type MedicalInformation struct {
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

func (medical *MedicalInformation) ToEntity() entities.MedicalInformation {
	return entities.MedicalInformation{
		Allergies:                  medical.Allergies,
		CurrentMedications:         medical.CurrentMedications,
		RecurrentMedicalConditions: medical.RecurrentMedicalConditions,
		HealthInsurancePlans:       medical.HealthInsurancePlans,
		BloodType:                  medical.BloodType,
		TakenVaccines:              medical.TakenVaccines,
		MentalHealthHistory:        medical.MentalHealthHistory,
		Height:                     medical.Height,
		Weight:                     medical.Weight,
		CigarettesUsage:            medical.CigarettesUsage,
		AlcoholConsumption:         medical.AlcoholConsumption,
		Disabilities:               medical.Disabilities,
	}
}
