package entities

import "time"

type BeneficiaryRoomAllocation struct {
	ID             string
	BeneficiaryID  string
	HousingID      string
	RoomID         string
	AllocationDate time.Time
	ExitDate       time.Time
	ExitReason     string
}
