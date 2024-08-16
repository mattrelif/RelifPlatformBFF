package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type Sessions interface {
	Generate(data entities.Session) (entities.Session, error)
	FindOneByUserID(userID string) (entities.Session, error)
	FindOneByIDAndUserID(id, userID string) (entities.Session, error)
	DeleteOneByID(id string) error
}

type sessionsMongo struct {
	collection *mongo.Collection
}

func NewSessionsMongo(database *mongo.Database) Sessions {
	return &sessionsMongo{
		collection: database.Collection("sessions"),
	}
}

func (repositories *sessionsMongo) Generate(data entities.Session) (entities.Session, error) {
	model := models.NewSession(data)

	if _, err := repositories.collection.InsertOne(context.Background(), model); err != nil {
		return entities.Session{}, err
	}

	return model.ToEntity(), nil
}

func (repositories *sessionsMongo) FindOneByUserID(userID string) (entities.Session, error) {
	var model models.Session

	filter := bson.M{"user_id": userID}

	if err := repositories.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Session{}, nil
		}

		return entities.Session{}, err
	}

	return model.ToEntity(), nil
}

func (repositories *sessionsMongo) FindOneByIDAndUserID(id, userID string) (entities.Session, error) {
	var model models.Session

	filter := bson.M{"$and": bson.A{
		bson.M{"user_id": userID},
		bson.M{"_id": id},
	},
	}

	if err := repositories.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Session{}, err
	}

	return model.ToEntity(), nil
}

func (repositories *sessionsMongo) DeleteOneByID(id string) error {
	filter := bson.M{"_id": id}
	if err := repositories.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
