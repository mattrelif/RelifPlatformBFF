package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
)

type BeneficiaryRoomAllocations interface {
	Create(allocation entities.BeneficiaryRoomAllocation) (string, error)
}

type mongoBeneficiaryRoomAllocations struct {
	collection *mongo.Collection
}

func NewMongoBeneficiaryRoomAllocations(database *mongo.Database) BeneficiaryRoomAllocations {
	return &mongoBeneficiaryRoomAllocations{
		collection: database.Collection("beneficiary_room_allocations"),
	}
}

func (repository *mongoBeneficiaryRoomAllocations) Create(allocation entities.BeneficiaryRoomAllocation) (string, error) {
	model := models.BeneficiaryRoomAllocation{
		ID:             primitive.NewObjectID().Hex(),
		BeneficiaryID:  allocation.BeneficiaryID,
		RoomID:         allocation.RoomID,
		AllocationDate: allocation.AllocationDate,
	}

	result, err := repository.collection.InsertOne(context.Background(), model)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
