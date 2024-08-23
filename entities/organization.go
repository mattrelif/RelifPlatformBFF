package entities

import "time"

type Organization struct {
	ID               string
	Name             string
	Description      string
	Address          Address
	Type             string
	OwnerID          string
	Status           string
	AccessGrantedIDs []string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
