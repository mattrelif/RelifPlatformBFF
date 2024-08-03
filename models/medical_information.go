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
	Addictions                 []string `bson:"addictions"`
	Disabilities               []string `bson:"disabilities"`
	ProthesisOrMedicalDevices  []string `bson:"prothesis_or_medical_devices"`
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
		Addictions:                 medical.Addictions,
		Disabilities:               medical.Disabilities,
		ProthesisOrMedicalDevices:  medical.ProthesisOrMedicalDevices,
	}
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
		Addictions:                 entity.Addictions,
		Disabilities:               entity.Disabilities,
		ProthesisOrMedicalDevices:  entity.ProthesisOrMedicalDevices,
	}
}
