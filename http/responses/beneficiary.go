package responses

import (
	"relif/bff/entities"
	"time"
)

type Beneficiaries []Beneficiary

type Beneficiary struct {
	ID                 string             `json:"id"`
	FullName           string             `json:"full_name"`
	Email              string             `json:"email"`
	Document           Document           `json:"document"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	CivilStatus        string             `json:"civil_status"`
	SpokenLanguages    []string           `json:"spoken_languages"`
	Education          string             `json:"education"`
	Address            Address            `json:"address"`
	CurrentHousingID   string             `json:"current_housing_id"`
	CurrentRoomID      string             `json:"current_room_id"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	Notes              []string           `json:"notes"`
}

func NewBeneficiary(entity entities.Beneficiary) Beneficiary {
	contactsEntityList := make([]EmergencyContact, 0)

	for _, contact := range entity.EmergencyContacts {
		contactsEntityList = append(contactsEntityList, NewEmergencyContact(contact))
	}

	return Beneficiary{
		ID:                 entity.ID,
		FullName:           entity.FullName,
		Email:              entity.Email,
		Document:           NewDocument(entity.Document),
		Birthdate:          entity.Birthdate,
		Phones:             entity.Phones,
		CivilStatus:        entity.CivilStatus,
		SpokenLanguages:    entity.SpokenLanguages,
		Education:          entity.Education,
		CurrentHousingID:   entity.CurrentHousingID,
		CurrentRoomID:      entity.CurrentRoomID,
		Address:            NewAddress(entity.Address),
		MedicalInformation: NewMedicalInformation(entity.MedicalInformation),
		EmergencyContacts:  contactsEntityList,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		Notes:              entity.Notes,
	}
}

func NewBeneficiaries(entityList []entities.Beneficiary) Beneficiaries {
	res := make(Beneficiaries, 0)

	for _, entity := range entityList {
		res = append(res, NewBeneficiary(entity))
	}

	return res
}
