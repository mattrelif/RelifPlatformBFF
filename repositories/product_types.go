package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
)

type ProductTypes interface {
	Create(data entities.ProductType) (entities.ProductType, error)
	FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.ProductType, error)
	FindOneByID(id string) (entities.ProductType, error)
	FindOneCompleteByID(id string) (entities.ProductType, error)
	UpdateOneByID(id string, data entities.ProductType) error
	DeleteOneByID(id string) error
}

type mongoProductTypes struct {
	collection *mongo.Collection
}

func NewMongoProductTypesRepository(database *mongo.Database) ProductTypes {
	return &mongoProductTypes{
		collection: database.Collection("product_types"),
	}
}

func (repository *mongoProductTypes) Create(data entities.ProductType) (entities.ProductType, error) {
	model := models.NewProductType(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.ProductType{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoProductTypes) FindOneByID(id string) (entities.ProductType, error) {
	var model models.ProductType

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.ProductType{}, utils.ErrProductTypeNotFound
		}

		return entities.ProductType{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoProductTypes) FindOneCompleteByID(id string) (entities.ProductType, error) {
	var model models.FindProductType

	filter := bson.M{"_id": id}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "organization_id"},
				{"foreignField", "_id"},
				{"as", "organization"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return entities.ProductType{}, err
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.ProductType{}, err
		}
	} else {
		return entities.ProductType{}, utils.ErrProductTypeNotFound
	}

	return model.ToEntity(), nil
}

func (repository *mongoProductTypes) FindManyByOrganizationID(organizationID string, limit, offset int64) (int64, []entities.ProductType, error) {
	modelList := make([]models.FindProductType, 0)
	entityList := make([]entities.ProductType, 0)

	filter := bson.M{"organization_id": organizationID}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$sort", bson.M{"name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "organization_id"},
				{"foreignField", "_id"},
				{"as", "organization"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

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

func (repository *mongoProductTypes) UpdateOneByID(id string, data entities.ProductType) error {
	model := models.NewUpdatedProductType(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoProductTypes) DeleteOneByID(id string) error {
	filter := bson.M{"_id": id}

	if _, err := repository.collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
