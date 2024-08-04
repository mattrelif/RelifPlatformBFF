package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type Housings interface {
	Create(data entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	UpdateOneById(id string, data entities.Housing) error
}

type mongoHousings struct {
	collection *mongo.Collection
}

func NewMongoHousings(database *mongo.Database) Housings {
	return &mongoHousings{
		collection: database.Collection("housings"),
	}
}

func (repository *mongoHousings) Create(data entities.Housing) (entities.Housing, error) {
	model := models.NewHousing(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	modelList := make([]models.Housing, 0)
	entityList := make([]entities.Housing, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"organization_id": organizationId,
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

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	defer cursor.Close(context.Background())

	if err != nil {
		return 0, nil, err
	}

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoHousings) FindOneByID(id string) (entities.Housing, error) {
	var model models.Housing

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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Housing{}, utils.ErrHousingNotFound
		}

		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) UpdateOneById(id string, data entities.Housing) error {
	model := models.NewUpdatedHousing(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
