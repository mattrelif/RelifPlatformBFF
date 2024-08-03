package entities

type MedicalInformation struct {
	Allergies                  []string
	CurrentMedications         []string
	RecurrentMedicalConditions []string
	HealthInsurancePlans       []string
	BloodType                  string
	TakenVaccines              []string
	MentalHealthHistory        []string
	Height                     int
	Weight                     int
	Addictions                 []string
	Disabilities               []string
	ProthesisOrMedicalDevices  []string
}
