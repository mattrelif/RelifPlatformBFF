package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type BeneficiaryAllocations []BeneficiaryAllocation

type BeneficiaryAllocation struct {
	ID            string    `json:"id"`
	BeneficiaryID string    `json:"beneficiary_id"`
	OldHousingID  string    `json:"old_housing_id"`
	OldRoomID     string    `json:"old_room_id"`
	HousingID     string    `json:"housing_id"`
	RoomID        string    `json:"room_id"`
	Type          string    `json:"type"`
	AuditorID     string    `json:"auditor_id"`
	Auditor       User      `json:"auditor,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	ExitDate      time.Time `json:"exit_date"`
	ExitReason    string    `json:"exit_reason"`
}

func NewBeneficiaryAllocation(entity entities.BeneficiaryAllocation) BeneficiaryAllocation {
	return BeneficiaryAllocation{
		ID:            entity.ID,
		BeneficiaryID: entity.BeneficiaryID,
		OldHousingID:  entity.OldHousingID,
		OldRoomID:     entity.OldRoomID,
		HousingID:     entity.HousingID,
		RoomID:        entity.RoomID,
		Type:          entity.Type,
		AuditorID:     entity.AuditorID,
		CreatedAt:     entity.CreatedAt,
		ExitDate:      entity.ExitDate,
		ExitReason:    entity.ExitReason,
	}
}

func NewBeneficiaryAllocations(entityList []entities.BeneficiaryAllocation) BeneficiaryAllocations {
	res := make(BeneficiaryAllocations, 0)

	for _, entity := range entityList {
		res = append(res, NewBeneficiaryAllocation(entity))
	}

	return res
}
