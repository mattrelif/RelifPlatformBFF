package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type BeneficiaryAllocations []BeneficiaryAllocation

type BeneficiaryAllocation struct {
	ID            string      `json:"id"`
	BeneficiaryID string      `json:"beneficiary_id"`
	OldHousingID  string      `json:"old_housing_id"`
	OldHousing    Housing     `json:"old_housing"`
	OldRoomID     string      `json:"old_room_id"`
	OldRoom       HousingRoom `json:"old_room"`
	HousingID     string      `json:"housing_id"`
	Housing       Housing     `json:"housing"`
	RoomID        string      `json:"room_id"`
	Room          HousingRoom `json:"room"`
	Type          string      `json:"type"`
	AuditorID     string      `json:"auditor_id"`
	Auditor       User        `json:"auditor,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	ExitDate      time.Time   `json:"exit_date"`
	ExitReason    string      `json:"exit_reason"`
}

func NewBeneficiaryAllocation(entity entities.BeneficiaryAllocation) BeneficiaryAllocation {
	return BeneficiaryAllocation{
		ID:            entity.ID,
		BeneficiaryID: entity.BeneficiaryID,
		OldHousingID:  entity.OldHousingID,
		OldHousing:    NewHousing(entity.OldHousing),
		OldRoomID:     entity.OldRoomID,
		OldRoom:       NewHousingRoom(entity.OldRoom),
		HousingID:     entity.HousingID,
		Housing:       NewHousing(entity.Housing),
		RoomID:        entity.RoomID,
		Room:          NewHousingRoom(entity.Room),
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
