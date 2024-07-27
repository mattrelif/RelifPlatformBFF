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

type Sessions interface {
	Generate(session entities.Session) (string, error)
	FindOneById(id string) (entities.Session, error)
	DeleteOneById(id string) error
}

type sessionsMongo struct {
	collection *mongo.Collection
}

func NewSessionsMongo(database *mongo.Database) Sessions {
	return &sessionsMongo{
		collection: database.Collection("sessions"),
	}
}

func (repositories *sessionsMongo) Generate(session entities.Session) (string, error) {
	model := models.Session{
		ID:        primitive.NewObjectID().Hex(),
		UserID:    session.UserID,
		ExpiresAt: session.ExpiresAt,
	}

	oid, err := primitive.ObjectIDFromHex(session.UserID)

	if err != nil {
		return "", err
	}

	filter := bson.M{"user_id": oid}
	update := bson.M{"$set": &model}
	opts := options.Update().SetUpsert(true)

	result, err := repositories.collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return "", err
	}

	return result.UpsertedID.(primitive.ObjectID).Hex(), nil
}

func (repositories *sessionsMongo) FindOneById(id string) (entities.Session, error) {
	var model models.Session

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.Session{}, err
	}

	filter := bson.M{"_id": oid}
	if err = repositories.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Session{}, err
	}

	return model.ToEntity(), nil
}

func (repositories *sessionsMongo) DeleteOneById(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err = repositories.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
