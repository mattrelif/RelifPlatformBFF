package entities

import "time"

type BeneficiaryAllocation struct {
	ID            string
	BeneficiaryID string
	OldHousingID  string
	OldHousing    Housing
	OldRoomID     string
	OldRoom       HousingRoom
	HousingID     string
	Housing       Housing
	RoomID        string
	Room          HousingRoom
	Type          string
	AuditorID     string
	CreatedAt     time.Time
	ExitDate      time.Time
	ExitReason    string
}
