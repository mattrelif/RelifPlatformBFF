package entities

import "time"

type Donation struct {
	ID             string
	OrganizationID string
	BeneficiaryID  string
	Beneficiary    Beneficiary
	From           Location
	ProductTypeID  string
	ProductType    ProductType
	Quantity       int
	CreatedAt      time.Time
}
