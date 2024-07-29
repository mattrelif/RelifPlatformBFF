package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
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

func (req *MedicalInformation) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(req.Allergies, validation.Each(validation.Required)),
		validation.Field(req.CurrentMedications, validation.Each(validation.Required)),
		validation.Field(req.RecurrentMedicalConditions, validation.Each(validation.Required)),
		validation.Field(req.HealthInsurancePlans, validation.Each(validation.Required)),
		validation.Field(req.BloodType, validation.Required),
		validation.Field(req.TakenVaccines, validation.Each(validation.Required)),
		validation.Field(req.MentalHealthHistory, validation.Each(validation.Required)),
		validation.Field(req.Weight, validation.Required),
		validation.Field(req.Height, validation.Required),
		validation.Field(req.CigarettesUsage, validation.Required),
		validation.Field(req.AlcoholConsumption, validation.Required),
		validation.Field(req.Disabilities, validation.Each(validation.Required)),
	)
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
		CigarettesUsage:            req.CigarettesUsage,
		AlcoholConsumption:         req.AlcoholConsumption,
		Disabilities:               req.Disabilities,
	}
}
