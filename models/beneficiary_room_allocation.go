package models

import (
	"relif/bff/entities"
	"time"
)

type BeneficiaryRoomAllocation struct {
	ID             string    `bson:"_id"`
	BeneficiaryID  string    `bson:"beneficiary_id"`
	HousingID      string    `bson:"housing_id"`
	RoomID         string    `bson:"room_id"`
	AllocationDate time.Time `bson:"allocation_date"`
	ExitDate       time.Time `bson:"exit_date"`
	ExitReason     string    `bson:"exit_reason"`
}

func (allocation *BeneficiaryRoomAllocation) ToEntity() entities.BeneficiaryRoomAllocation {
	return entities.BeneficiaryRoomAllocation{
		ID:             allocation.ID,
		BeneficiaryID:  allocation.BeneficiaryID,
		RoomID:         allocation.RoomID,
		AllocationDate: allocation.AllocationDate,
		ExitDate:       allocation.ExitDate,
		ExitReason:     allocation.ExitReason,
	}
}
