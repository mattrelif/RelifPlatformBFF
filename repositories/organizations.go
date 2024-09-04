package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
)

type Organizations interface {
	Create(data entities.Organization) (entities.Organization, error)
	FindManyPaginated(offset, limit int64) (int64, []entities.Organization, error)
	FindOneByID(id string) (entities.Organization, error)
	UpdateOneByID(id string, data entities.Organization) error
}

type mongoOrganizations struct {
	collection *mongo.Collection
}

func NewMongoOrganizations(database *mongo.Database) Organizations {
	return &mongoOrganizations{
		collection: database.Collection("organizations"),
	}
}

func (repository *mongoOrganizations) Create(data entities.Organization) (entities.Organization, error) {
	model := models.NewOrganization(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.Organization{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizations) FindManyPaginated(offset, limit int64) (int64, []entities.Organization, error) {
	entityList := make([]entities.Organization, 0)

	count, err := repository.collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
	cursor, err := repository.collection.Find(context.Background(), bson.M{}, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.Organization

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizations) FindOneByID(id string) (entities.Organization, error) {
	var model models.Organization

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Organization{}, utils.ErrOrganizationNotFound
		}
		return entities.Organization{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizations) UpdateOneByID(id string, data entities.Organization) error {
	model := models.NewUpdatedOrganization(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
