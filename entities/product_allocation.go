package entities

import "time"

type ProductAllocation struct {
	ID             string
	ProductTypeID  string
	ProductType    ProductType
	Type           string
	FromHousingID  string
	FromHousing    Housing
	ToHousingID    string
	ToHousing      Housing
	OrganizationID string
	Organization   Organization
	Quantity       int
	AuditorID      string
	Auditor        User
	CreatedAt      time.Time
}
