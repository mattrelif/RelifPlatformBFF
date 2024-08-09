package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type UpdateHousingRoom struct {
	Name           string `json:"name"`
	Status         string `json:"status"`
	TotalVacancies int    `json:"total_vacancies"`
}

func (req *UpdateHousingRoom) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.Status, validation.Required),
		validation.Field(&req.TotalVacancies, validation.Required),
	)
}

func (req *UpdateHousingRoom) ToEntity() entities.HousingRoom {
	return entities.HousingRoom{
		Name:           req.Name,
		Status:         req.Status,
		TotalVacancies: req.TotalVacancies,
	}
}
