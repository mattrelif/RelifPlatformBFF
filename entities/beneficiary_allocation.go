package entities

import "time"

type BeneficiaryAllocation struct {
	ID            string
	BeneficiaryID string
	OldHousingID  string
	OldRoomID     string
	HousingID     string
	RoomID        string
	Type          string
	AuditorID     string
	CreatedAt     time.Time
	ExitDate      time.Time
	ExitReason    string
}
