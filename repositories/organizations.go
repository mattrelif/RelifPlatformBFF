package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type Organizations interface {
	Create(data entities.Organization) (entities.Organization, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
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

func (repository *mongoOrganizations) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	modelList := make([]models.Organization, 0)
	entityList := make([]entities.Organization, 0)

	filter := bson.M{"status": bson.M{"$not": bson.M{"$eq": utils.InactiveStatus}}}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
	cursor, err := repository.collection.Find(context.Background(), bson.M{}, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizations) FindOneByID(id string) (entities.Organization, error) {
	var model models.Organization

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
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
