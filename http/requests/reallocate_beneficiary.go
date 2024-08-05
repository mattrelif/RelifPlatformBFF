package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"relif/bff/entities"
)

type ReallocateBeneficiary struct {
	HousingID  string `json:"housing_id"`
	RoomID     string `json:"room_id"`
	ExitReason string `json:"exit_reason"`
}

func (req *ReallocateBeneficiary) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.HousingID, validation.Required, is.MongoID),
		validation.Field(&req.ExitReason, validation.Required),
	)
}

func (req *ReallocateBeneficiary) ToEntity() entities.BeneficiaryAllocation {
	return entities.BeneficiaryAllocation{
		HousingID:  req.HousingID,
		RoomID:     req.RoomID,
		ExitReason: req.ExitReason,
	}
}
