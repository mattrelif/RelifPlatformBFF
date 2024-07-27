package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
)

type OrganizationDataAccessGrants interface {
	Create(grant entities.OrganizationDataAccessGrant) error
}

type mongoOrganizationDataAccessGrants struct {
	collection *mongo.Collection
}

func NewMongoOrganizationDataAccessGrants(database *mongo.Database) OrganizationDataAccessGrants {
	return &mongoOrganizationDataAccessGrants{
		collection: database.Collection("organization_data_access_grants"),
	}
}

func (m *mongoOrganizationDataAccessGrants) Create(access entities.OrganizationDataAccessGrant) error {
	model := models.OrganizationDataAccess{
		ID:                   primitive.NewObjectID().Hex(),
		OrganizationID:       access.OrganizationID,
		TargetOrganizationID: access.TargetOrganizationID,
		AuditorID:            access.AuditorID,
		CreatedAt:            access.CreatedAt,
	}

	if _, err := m.collection.InsertOne(context.Background(), &model); err != nil {
		return err
	}

	return nil
}
