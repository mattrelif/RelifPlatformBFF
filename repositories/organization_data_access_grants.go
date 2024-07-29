package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
	"time"
)

type OrganizationDataAccessGrants interface {
	Create(data entities.OrganizationDataAccessGrant) error
}

type mongoOrganizationDataAccessGrants struct {
	collection *mongo.Collection
}

func NewMongoOrganizationDataAccessGrants(database *mongo.Database) OrganizationDataAccessGrants {
	return &mongoOrganizationDataAccessGrants{
		collection: database.Collection("organization_data_access_grants"),
	}
}

func (m *mongoOrganizationDataAccessGrants) Create(data entities.OrganizationDataAccessGrant) error {
	model := models.OrganizationDataAccess{
		ID:                   primitive.NewObjectID().Hex(),
		OrganizationID:       data.OrganizationID,
		TargetOrganizationID: data.TargetOrganizationID,
		AuditorID:            data.AuditorID,
		CreatedAt:            time.Now(),
	}

	if _, err := m.collection.InsertOne(context.Background(), &model); err != nil {
		return err
	}

	return nil
}
