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

type Housings interface {
	Create(housing entities.Housing) (string, error)
	FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	FindOneAndUpdateById(id string, housing entities.Housing) (entities.Housing, error)
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

func (repository *mongoHousings) Create(housing entities.Housing) (string, error) {
	model := models.Housing{
		ID:             primitive.NewObjectID().Hex(),
		Name:           housing.Name,
		OrganizationID: housing.OrganizationID,
		Status:         housing.Status,
		Address: models.Address{
			StreetName:   housing.Address.StreetName,
			StreetNumber: housing.Address.StreetNumber,
			City:         housing.Address.City,
			ZipCode:      housing.Address.ZipCode,
			Country:      housing.Address.Country,
			District:     housing.Address.District,
		},
		CreatedAt: housing.CreatedAt,
	}

	res, err := repository.collection.InsertOne(context.Background(), model)

	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repository *mongoHousings) FindManyByOrganizationID(organizationId string, limit, offset int64) (int64, []entities.Housing, error) {
	modelList := make([]models.Housing, 0)
	entityList := make([]entities.Housing, 0)

	oid, err := primitive.ObjectIDFromHex(organizationId)

	if err != nil {
		return 0, nil, err
	}

	filter := bson.M{"organization_id": oid}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

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

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.Housing{}, err
	}

	filter := bson.M{"_id": oid}

	if err = repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindOneAndUpdateById(id string, housing entities.Housing) (entities.Housing, error) {
	model := models.Housing{
		Name:   housing.Name,
		Status: housing.Status,
		Address: models.Address{
			StreetName:   housing.Address.StreetName,
			StreetNumber: housing.Address.StreetNumber,
			City:         housing.Address.City,
			ZipCode:      housing.Address.ZipCode,
			Country:      housing.Address.Country,
			District:     housing.Address.District,
		},
		UpdatedAt: housing.UpdatedAt,
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.Housing{}, err
	}

	filter := bson.M{"_id": oid}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, model, opts).Decode(&model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) DeleteOneById(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	if err = repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
