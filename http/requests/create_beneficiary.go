package requests

import (
	"relif/platform-bff/entities"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateBeneficiary struct {
	FullName           string             `json:"full_name"`
	Email              string             `json:"email"`
	ImageURL           string             `json:"image_url"`
	Documents          []Document         `json:"documents"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	CivilStatus        string             `json:"civil_status"`
	SpokenLanguages    []string           `json:"spoken_languages"`
	Education          string             `json:"education"`
	Occupation         string             `json:"occupation"`
	Address            Address            `json:"address"`
	Gender             string             `json:"gender"`
	Status             string             `json:"status,omitempty"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	Notes              string             `json:"notes"`
}

func (req *CreateBeneficiary) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.FullName, validation.Required),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Documents, validation.Each(validation.By(func(value interface{}) error {
			if document, ok := value.(Document); ok {
				return document.Validate()
			}

			return nil
		}))),
		validation.Field(&req.Birthdate, validation.Required),
		validation.Field(&req.Phones, validation.Each(validation.Required)),
		validation.Field(&req.CivilStatus, validation.Required),
		validation.Field(&req.Education, validation.Required),
		validation.Field(&req.SpokenLanguages, validation.Each(validation.Required)),
		validation.Field(&req.Gender, validation.Required),
		validation.Field(&req.Occupation, validation.Required),
		validation.Field(&req.Status, validation.When(req.Status != "", validation.In("ACTIVE", "INACTIVE", "PENDING", "ARCHIVED"))),
		validation.Field(&req.Address, validation.By(func(value interface{}) error {
			if address, ok := value.(Address); ok {
				return address.Validate()
			}
			return nil
		})),
		validation.Field(&req.EmergencyContacts, validation.By(func(value interface{}) error {
			if contacts, ok := value.([]EmergencyContact); ok {
				for _, contact := range contacts {
					if err := contact.Validate(); err != nil {
						return err
					}
				}
			}
			return nil
		})),
	)
}

func (req *CreateBeneficiary) ToEntity() entities.Beneficiary {
	contactsEntityList := make([]entities.EmergencyContact, 0)

	for _, contact := range req.EmergencyContacts {
		contactsEntityList = append(contactsEntityList, contact.ToEntity())
	}

	documentsEntityList := make([]entities.Document, 0)

	for _, document := range req.Documents {
		documentsEntityList = append(documentsEntityList, document.ToEntity())
	}

	return entities.Beneficiary{
		FullName:           req.FullName,
		Email:              req.Email,
		ImageURL:           req.ImageURL,
		Documents:          documentsEntityList,
		Birthdate:          req.Birthdate,
		Phones:             req.Phones,
		CivilStatus:        req.CivilStatus,
		SpokenLanguages:    req.SpokenLanguages,
		Education:          req.Education,
		Gender:             req.Gender,
		Occupation:         req.Occupation,
		Address:            req.Address.ToEntity(),
		Status:             req.Status,
		MedicalInformation: req.MedicalInformation.ToEntity(),
		EmergencyContacts:  contactsEntityList,
		Notes:              req.Notes,
	}
}
