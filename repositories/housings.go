package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"time"
)

type Housings interface {
	Create(data entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	FindOneAndUpdateById(id string, data entities.Housing) (entities.Housing, error)
	DeleteOneById(id string) error
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
	model := models.Housing{
		ID:             primitive.NewObjectID().Hex(),
		Name:           data.Name,
		OrganizationID: data.OrganizationID,
		Status:         data.Status,
		Address: models.Address{
			StreetName:   data.Address.StreetName,
			StreetNumber: data.Address.StreetNumber,
			City:         data.Address.City,
			ZipCode:      data.Address.ZipCode,
			Country:      data.Address.Country,
			District:     data.Address.District,
		},
		CreatedAt: time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	modelList := make([]models.Housing, 0)
	entityList := make([]entities.Housing, 0)

	filter := bson.M{"organization_id": organizationId}
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

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindOneAndUpdateById(id string, data entities.Housing) (entities.Housing, error) {
	model := models.Housing{
		Name:   data.Name,
		Status: data.Status,
		Address: models.Address{
			StreetName:   data.Address.StreetName,
			StreetNumber: data.Address.StreetNumber,
			City:         data.Address.City,
			ZipCode:      data.Address.ZipCode,
			Country:      data.Address.Country,
			District:     data.Address.District,
		},
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, model, opts).Decode(&model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
