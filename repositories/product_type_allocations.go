package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type ProductTypeAllocations interface {
	Create(data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
}

type mongoProductTypeAllocations struct {
	collection *mongo.Collection
}

func NewMongoProductTypeAllocations(database *mongo.Database) ProductTypeAllocations {
	return &mongoProductTypeAllocations{
		collection: database.Collection("product_type_allocations"),
	}
}

func (repository *mongoProductTypeAllocations) Create(data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	model := models.NewProductTypeAllocation(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	return data, nil
}
