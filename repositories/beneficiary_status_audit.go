package repositories

import (
	"context"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type BeneficiaryStatusAudit interface {
	Create(data entities.BeneficiaryStatusAudit) error
	FindManyByBeneficiaryID(beneficiaryID string) ([]entities.BeneficiaryStatusAudit, error)
}

type mongoBeneficiaryStatusAudit struct {
	collection *mongo.Collection
}

func NewMongoBeneficiaryStatusAudit(database *mongo.Database) BeneficiaryStatusAudit {
	return &mongoBeneficiaryStatusAudit{
		collection: database.Collection("beneficiary_status_audit"),
	}
}

func (repository *mongoBeneficiaryStatusAudit) Create(data entities.BeneficiaryStatusAudit) error {
	model := models.NewBeneficiaryStatusAudit(data)

	_, err := repository.collection.InsertOne(context.Background(), &model)

	return err
}

func (repository *mongoBeneficiaryStatusAudit) FindManyByBeneficiaryID(beneficiaryID string) ([]entities.BeneficiaryStatusAudit, error) {
	// Implementation for finding audit records by beneficiary ID (for future use)
	return []entities.BeneficiaryStatusAudit{}, nil
}
