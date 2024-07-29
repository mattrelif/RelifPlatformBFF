package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type AllocateBeneficiary struct {
	HousingID string `json:"housing_id"`
	RoomID    string `json:"room_id"`
}

func (req *AllocateBeneficiary) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.HousingID, validation.Required, is.MongoID),
		validation.Field(&req.RoomID, validation.Required, is.MongoID),
	)
}

func (req *AllocateBeneficiary) ToEntity() entities.BeneficiaryAllocation {
	return entities.BeneficiaryAllocation{
		HousingID: req.HousingID,
		RoomID:    req.RoomID,
	}
}
