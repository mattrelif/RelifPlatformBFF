package requests

import (
	"relif/platform-bff/entities"
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
	Addictions                 []string `json:"addictions"`
	Disabilities               []string `json:"disabilities"`
	ProthesisOrMedicalDevices  []string `json:"prothesis_or_medical_devices"`
}

func (req *MedicalInformation) ToEntity() entities.MedicalInformation {
	return entities.MedicalInformation{
		Allergies:                  req.Allergies,
		CurrentMedications:         req.CurrentMedications,
		RecurrentMedicalConditions: req.RecurrentMedicalConditions,
		HealthInsurancePlans:       req.HealthInsurancePlans,
		BloodType:                  req.BloodType,
		TakenVaccines:              req.TakenVaccines,
		MentalHealthHistory:        req.MentalHealthHistory,
		Height:                     req.Height,
		Weight:                     req.Weight,
		Addictions:                 req.Addictions,
		Disabilities:               req.Disabilities,
		ProthesisOrMedicalDevices:  req.ProthesisOrMedicalDevices,
	}
}
