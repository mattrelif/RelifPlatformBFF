package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"relif/platform-bff/entities"
)

type CreateManyHousingRooms []CreateHousingRoom

type CreateHousingRoom struct {
	Name           string `json:"name"`
	TotalVacancies int    `json:"total_vacancies"`
}

func (req *CreateHousingRoom) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required),
		validation.Field(&req.TotalVacancies, validation.Required),
	)
}

func (req *CreateHousingRoom) ToEntity() entities.HousingRoom {
	return entities.HousingRoom{
		Name:           req.Name,
		TotalVacancies: req.TotalVacancies,
	}
}

func (req *CreateManyHousingRooms) Validate() error {
	for _, r := range *req {
		if err := r.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (req *CreateManyHousingRooms) ToEntity() []entities.HousingRoom {
	entityList := make([]entities.HousingRoom, 0)

	for _, r := range *req {
		entityList = append(entityList, r.ToEntity())
	}

	return entityList
}
