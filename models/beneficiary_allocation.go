package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"time"
)

type BeneficiaryAllocation struct {
	ID            string    `bson:"_id,omitempty"`
	BeneficiaryID string    `bson:"beneficiary_id,omitempty"`
	OldHousingID  string    `bson:"old_housing_id,omitempty"`
	OldRoomID     string    `bson:"old_room_id,omitempty"`
	HousingID     string    `bson:"housing_id,omitempty"`
	RoomID        string    `bson:"room_id,omitempty"`
	Type          string    `bson:"type,omitempty"`
	AuditorID     string    `bson:"auditor_id,omitempty"`
	CreatedAt     time.Time `bson:"created_at,omitempty"`
	ExitDate      time.Time `bson:"exit_date,omitempty"`
	ExitReason    string    `bson:"exit_reason,omitempty"`
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

func NewBeneficiaryAllocation(entity entities.BeneficiaryAllocation) BeneficiaryAllocation {
	return BeneficiaryAllocation{
		ID:            primitive.NewObjectID().Hex(),
		BeneficiaryID: entity.BeneficiaryID,
		OldHousingID:  entity.OldHousingID,
		OldRoomID:     entity.OldRoomID,
		HousingID:     entity.HousingID,
		RoomID:        entity.RoomID,
		Type:          entity.Type,
		AuditorID:     entity.AuditorID,
		CreatedAt:     time.Now(),
		ExitDate:      entity.ExitDate,
		ExitReason:    entity.ExitReason,
	}
}
