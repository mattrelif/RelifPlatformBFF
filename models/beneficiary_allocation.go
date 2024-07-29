package models

import (
	"relif/bff/entities"
	"time"
)

type BeneficiaryAllocation struct {
	ID            string    `bson:"_id"`
	BeneficiaryID string    `bson:"beneficiary_id"`
	OldHousingID  string    `bson:"old_housing_id"`
	OldRoomID     string    `bson:"old_room_id"`
	HousingID     string    `bson:"housing_id"`
	RoomID        string    `bson:"room_id"`
	Type          string    `bson:"type"`
	AuditorID     string    `bson:"auditor_id"`
	CreatedAt     time.Time `bson:"created_at"`
	ExitDate      time.Time `bson:"exit_date"`
	ExitReason    string    `bson:"exit_reason"`
}

func (allocation *BeneficiaryAllocation) ToEntity() entities.BeneficiaryAllocation {
	return entities.BeneficiaryAllocation{
		ID:            allocation.ID,
		BeneficiaryID: allocation.BeneficiaryID,
		OldHousingID:  allocation.OldHousingID,
		OldRoomID:     allocation.OldRoomID,
		HousingID:     allocation.HousingID,
		RoomID:        allocation.RoomID,
		Type:          allocation.Type,
		AuditorID:     allocation.AuditorID,
		CreatedAt:     allocation.CreatedAt,
		ExitDate:      allocation.ExitDate,
		ExitReason:    allocation.ExitReason,
	}
}
