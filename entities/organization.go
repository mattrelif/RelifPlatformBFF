package entities

import "time"

type Organization struct {
	ID          string
	Name        string
	Description string
	Address     Address
	Type        string
	CreatorID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
