package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
)

type PasswordChangeRequests interface {
	Create(request entities.PasswordChangeRequest) (string, error)
	FindOneAndDeleteById(id string) (entities.PasswordChangeRequest, error)
}

type mongoPasswordChangeRequests struct {
	collection *mongo.Collection
}

func NewMongoPasswordChangeRequests(database *mongo.Database) PasswordChangeRequests {
	return &mongoPasswordChangeRequests{
		collection: database.Collection("password_change_requests"),
	}
}

func (rep *mongoPasswordChangeRequests) Create(request entities.PasswordChangeRequest) (string, error) {
	model := models.PasswordChangeRequest{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    request.UserID,
		ExpiresAt: request.ExpiresAt,
	}

	oid, err := primitive.ObjectIDFromHex(request.UserID)

	if err != nil {
		return "", err
	}

	filter := bson.M{"user_id": oid}
	update := bson.M{"$set": &model}
	opts := options.Update().SetUpsert(true)

	result, err := rep.collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return "", err
	}

	return result.UpsertedID.(primitive.ObjectID).Hex(), nil
}

func (rep *mongoPasswordChangeRequests) FindOneAndDeleteById(id string) (entities.PasswordChangeRequest, error) {
	var model models.PasswordChangeRequest

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.PasswordChangeRequest{}, err
	}

	filter := bson.M{"user_id": oid}
	if err = rep.collection.FindOneAndDelete(context.Background(), filter).Decode(&model); err != nil {
		return entities.PasswordChangeRequest{}, err
	}

	return model.ToEntity(), nil
}
