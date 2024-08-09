package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/platform-bff/entities"
)

type CreateVoluntaryPerson struct {
	FullName           string             `json:"full_name"`
	Email              string             `json:"email"`
	Documents          []Document         `json:"documents"`
	Birthdate          string             `json:"birthdate"`
	Phones             []string           `json:"phones"`
	Segments           []string           `json:"segments"`
	Gender             string             `json:"gender"`
	Address            Address            `json:"address"`
	MedicalInformation MedicalInformation `json:"medical_information"`
	EmergencyContacts  []EmergencyContact `json:"emergency_contacts"`
	Notes              string             `json:"notes"`
}

func (req *CreateVoluntaryPerson) Validate() error {
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
		validation.Field(&req.Segments, validation.Each(validation.Required)),
		validation.Field(&req.Phones, validation.Each(validation.Required)),
		validation.Field(&req.Gender, validation.Required),
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

func (req *CreateVoluntaryPerson) ToEntity() entities.VoluntaryPerson {
	contactsEntityList := make([]entities.EmergencyContact, 0)

	for _, contact := range req.EmergencyContacts {
		contactsEntityList = append(contactsEntityList, contact.ToEntity())
	}

	documentsEntityList := make([]entities.Document, 0)

	for _, document := range req.Documents {
		documentsEntityList = append(documentsEntityList, document.ToEntity())
	}

	return entities.VoluntaryPerson{
		FullName:           req.FullName,
		Email:              req.Email,
		Documents:          documentsEntityList,
		Birthdate:          req.Birthdate,
		Phones:             req.Phones,
		Address:            req.Address.ToEntity(),
		Segments:           req.Segments,
		Gender:             req.Gender,
		MedicalInformation: req.MedicalInformation.ToEntity(),
		EmergencyContacts:  contactsEntityList,
		Notes:              req.Notes,
	}
}
