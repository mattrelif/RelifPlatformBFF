package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type Donations []Donation

type Donation struct {
	ID             string      `json:"id,omitempty"`
	OrganizationID string      `json:"organization_id,omitempty"`
	BeneficiaryID  string      `json:"beneficiary_id,omitempty"`
	Beneficiary    Beneficiary `json:"beneficiary,omitempty"`
	From           Location    `json:"from,omitempty"`
	ProductTypeID  string      `json:"product_type_id,omitempty"`
	ProductType    ProductType `json:"product_type,omitempty"`
	Quantity       int         `json:"quantity,omitempty"`
	CreatedAt      time.Time   `json:"created_at,omitempty"`
}

func NewDonation(entity entities.Donation) Donation {
	return Donation{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		BeneficiaryID:  entity.BeneficiaryID,
		Beneficiary:    NewBeneficiary(entity.Beneficiary),
		From:           NewLocation(entity.From),
		ProductTypeID:  entity.ProductTypeID,
		ProductType:    NewProductType(entity.ProductType),
		Quantity:       entity.Quantity,
		CreatedAt:      entity.CreatedAt,
	}
}

func NewDonations(entityList []entities.Donation) Donations {
	res := make(Donations, 0)

	for _, entity := range entityList {
		res = append(res, NewDonation(entity))
	}

	return res
}
