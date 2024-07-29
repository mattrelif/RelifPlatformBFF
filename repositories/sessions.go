package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
)

type Sessions interface {
	Generate(data entities.Session) error
	FindOneBySessionId(sessionId string) (entities.Session, error)
	DeleteOneBySessionId(sessionId string) error
}

type sessionsMongo struct {
	collection *mongo.Collection
}

func NewSessionsMongo(database *mongo.Database) Sessions {
	return &sessionsMongo{
		collection: database.Collection("sessions"),
	}
}

func (repositories *sessionsMongo) Generate(data entities.Session) error {
	model := models.Session{
		UserID:    data.UserID,
		SessionID: data.SessionID,
		ExpiresAt: data.ExpiresAt,
	}

	filter := bson.M{"_id": model.UserID}
	update := bson.M{"$set": &model}
	opts := options.Update().SetUpsert(true)

	if _, err := repositories.collection.UpdateOne(context.Background(), filter, update, opts); err != nil {
		return err
	}

	return nil
}

func (repositories *sessionsMongo) FindOneBySessionId(sessionId string) (entities.Session, error) {
	var model models.Session

	filter := bson.M{"session_id": sessionId}
	if err := repositories.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Session{}, err
	}

	return model.ToEntity(), nil
}

func (repositories *sessionsMongo) DeleteOneBySessionId(sessionId string) error {
	filter := bson.M{"session_id": sessionId}
	if err := repositories.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
