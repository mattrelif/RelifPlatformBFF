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

type Organizations interface {
	Create(data entities.Organization) (entities.Organization, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneAndUpdateById(id string, data entities.Organization) (entities.Organization, error)
	UpdateOneById(id string, data entities.Organization) error
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
	model := models.Organization{
		ID:          primitive.NewObjectID().Hex(),
		Name:        data.Name,
		Description: data.Description,
		Address: models.Address{
			StreetName:   data.Address.StreetName,
			StreetNumber: data.Address.StreetNumber,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Type:      data.Type,
		CreatorID: data.CreatorID,
		CreatedAt: time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.Organization{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizations) FindMany(offset, limit int64) (int64, []entities.Organization, error) {
	modelList := make([]models.Organization, 0)
	entityList := make([]entities.Organization, 0)

	count, err := repository.collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
	cursor, err := repository.collection.Find(context.Background(), bson.M{}, opts)

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

func (repository *mongoOrganizations) FindOneAndUpdateById(id string, data entities.Organization) (entities.Organization, error) {
	model := models.Organization{
		Name:        data.Name,
		Description: data.Description,
		Address: models.Address{
			StreetName:   data.Address.StreetName,
			StreetNumber: data.Address.StreetNumber,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Type:      data.Type,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.Organization{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizations) UpdateOneById(id string, data entities.Organization) error {
	model := models.Organization{
		Name:        data.Name,
		Description: data.Description,
		Address: models.Address{
			StreetName:   data.Address.StreetName,
			StreetNumber: data.Address.StreetNumber,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Type:      data.Type,
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Err(); err != nil {
		return err
	}

	return nil
}
