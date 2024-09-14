package entities

import "time"

type ProductTypeAllocation struct {
	ID             string
	ProductTypeID  string
	ProductType    ProductType
	Type           string
	From           Location
	To             Location
	OrganizationID string
	Quantity       int
	CreatedAt      time.Time
}
