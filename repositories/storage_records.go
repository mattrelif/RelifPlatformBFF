package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type StorageRecords interface {
	Create(data entities.StorageRecord) error
	FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error)
	UpdateOneByID(id string, data entities.StorageRecord) error
}

type mongoStorageRecords struct {
	collection *mongo.Collection
}

func NewMongoStorageRecords(database *mongo.Database) StorageRecords {
	return &mongoStorageRecords{
		collection: database.Collection("storage_records"),
	}
}

func (repository *mongoStorageRecords) Create(data entities.StorageRecord) error {
	model := models.NewStorageRecord(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return err
	}

	return nil
}

func (repository *mongoStorageRecords) FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error) {
	var model models.StorageRecord

	filter := bson.M{
		"product_type_id": productTypeID,
		"location.id":     location.ID,
		"location.type":   location.Type,
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.StorageRecord{}, nil
		}

		return entities.StorageRecord{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoStorageRecords) UpdateOneByID(id string, data entities.StorageRecord) error {
	model := models.NewUpdatedStorageRecord(data)

	update := bson.M{
		"$set": model,
	}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
