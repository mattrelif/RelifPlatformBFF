package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"time"
)

type FindDonation struct {
	ID             string       `bson:"_id,omitempty"`
	OrganizationID string       `bson:"organization_id,omitempty"`
	BeneficiaryID  string       `bson:"beneficiary_id,omitempty"`
	Beneficiary    Beneficiary  `bson:"beneficiary,omitempty"`
	From           FindLocation `bson:"from,omitempty"`
	ProductTypeID  string       `bson:"product_type_id,omitempty"`
	ProductType    ProductType  `bson:"product_type,omitempty"`
	Quantity       int          `bson:"quantity,omitempty"`
	CreatedAt      time.Time    `bson:"created_at,omitempty"`
}

func (donation *FindDonation) ToEntity() entities.Donation {
	return entities.Donation{
		ID:             donation.ID,
		OrganizationID: donation.OrganizationID,
		BeneficiaryID:  donation.BeneficiaryID,
		Beneficiary:    donation.Beneficiary.ToEntity(),
		From:           donation.From.ToEntity(),
		ProductTypeID:  donation.ProductTypeID,
		ProductType:    donation.ProductType.ToEntity(),
		Quantity:       donation.Quantity,
		CreatedAt:      donation.CreatedAt,
	}
}

type Donation struct {
	ID             string    `bson:"_id,omitempty"`
	OrganizationID string    `bson:"organization_id,omitempty"`
	BeneficiaryID  string    `bson:"beneficiary_id,omitempty"`
	From           Location  `bson:"from,omitempty"`
	ProductTypeID  string    `bson:"product_type_id,omitempty"`
	Quantity       int       `bson:"quantity,omitempty"`
	CreatedAt      time.Time `bson:"created_at,omitempty"`
}

func (donation *Donation) ToEntity() entities.Donation {
	return entities.Donation{
		ID:             donation.ID,
		OrganizationID: donation.OrganizationID,
		BeneficiaryID:  donation.BeneficiaryID,
		From:           donation.From.ToEntity(),
		ProductTypeID:  donation.ProductTypeID,
		Quantity:       donation.Quantity,
		CreatedAt:      donation.CreatedAt,
	}
}

func NewDonation(entity entities.Donation) Donation {
	return Donation{
		ID:             primitive.NewObjectID().Hex(),
		OrganizationID: entity.OrganizationID,
		BeneficiaryID:  entity.BeneficiaryID,
		From:           NewLocation(entity.From),
		ProductTypeID:  entity.ProductTypeID,
		Quantity:       entity.Quantity,
		CreatedAt:      time.Now(),
	}
}
