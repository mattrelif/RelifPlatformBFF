package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
	"time"
)

type ProductTypes interface {
	Create(data entities.ProductType) (entities.ProductType, error)
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.ProductType, error)
	FindOneById(id string) (entities.ProductType, error)
	UpdateOneById(id string, data entities.ProductType) error
	IncreaseTotalInStock(id string, amount int) error
	DeleteOneById(id string) error
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
	model := models.ProductType{
		ID:             primitive.NewObjectID().Hex(),
		Name:           data.Name,
		Description:    data.Description,
		Brand:          data.Brand,
		Category:       data.Category,
		OrganizationID: data.OrganizationID,
		CreatedAt:      time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.TODO(), model); err != nil {
		return entities.ProductType{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoProductTypes) FindOneById(id string) (entities.ProductType, error) {
	var model models.ProductType

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.TODO(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.ProductType{}, utils.ErrProductTypeNotFound
		}
		return entities.ProductType{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoProductTypes) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.ProductType, error) {
	modelList := make([]models.ProductType, 0)
	entityList := make([]entities.ProductType, 0)

	filter := bson.M{"organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})

	cursor, err := repository.collection.Find(context.TODO(), filter, opts)

	if err != nil {
		return 0, nil, err
	}

	if err = cursor.All(context.TODO(), &modelList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoProductTypes) UpdateOneById(id string, data entities.ProductType) error {
	model := models.NewUpdatedProductType(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.TODO(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoProductTypes) IncreaseTotalInStock(id string, amount int) error {
	update := bson.M{"$inc": bson.M{"total_in_stock": amount}, "$set": bson.M{"updated_at": time.Now()}}

	if _, err := repository.collection.UpdateByID(context.TODO(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoProductTypes) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.TODO(), filter).Err(); err != nil {
		return err
	}

	return nil
}
