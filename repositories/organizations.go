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

type Organizations interface {
	Create(organization entities.Organization) (string, error)
	FindMany(offset, limit int64) (int64, []entities.Organization, error)
	FindOneAndUpdateById(id string, organization entities.Organization) (entities.Organization, error)
	UpdateOneById(id string, organization entities.Organization) error
}

type mongoOrganizations struct {
	collection *mongo.Collection
}

func NewMongoOrganizations(database *mongo.Database) Organizations {
	return &mongoOrganizations{
		collection: database.Collection("organizations"),
	}
}

func (repository *mongoOrganizations) Create(organization entities.Organization) (string, error) {
	model := models.Organization{
		ID:          primitive.NewObjectID().Hex(),
		Name:        organization.Name,
		Description: organization.Description,
		Address: models.OrganizationAddress{
			StreetName:   organization.Address.StreetName,
			StreetNumber: organization.Address.StreetNumber,
			ZipCode:      organization.Address.ZipCode,
			District:     organization.Address.District,
			City:         organization.Address.City,
			Country:      organization.Address.Country,
		},
		Type:      organization.Type,
		CreatorID: organization.CreatorID,
		CreatedAt: organization.CreatedAt,
	}

	result, err := repository.collection.InsertOne(context.Background(), &model)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
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

func (repository *mongoOrganizations) FindOneAndUpdateById(id string, organization entities.Organization) (entities.Organization, error) {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.Organization{}, err
	}

	model := models.Organization{
		Name:        organization.Name,
		Description: organization.Description,
		Address: models.OrganizationAddress{
			StreetName:   organization.Address.StreetName,
			StreetNumber: organization.Address.StreetNumber,
			ZipCode:      organization.Address.ZipCode,
			District:     organization.Address.District,
			City:         organization.Address.City,
			Country:      organization.Address.Country,
		},
		Type:      organization.Type,
		CreatorID: organization.CreatorID,
		UpdatedAt: organization.UpdatedAt,
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.Organization{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizations) UpdateOneById(id string, organization entities.Organization) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	model := models.Organization{
		Name:        organization.Name,
		Description: organization.Description,
		Address: models.OrganizationAddress{
			StreetName:   organization.Address.StreetName,
			StreetNumber: organization.Address.StreetNumber,
			ZipCode:      organization.Address.ZipCode,
			District:     organization.Address.District,
			City:         organization.Address.City,
			Country:      organization.Address.Country,
		},
		Type:      organization.Type,
		CreatorID: organization.CreatorID,
		UpdatedAt: organization.UpdatedAt,
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Err(); err != nil {
		return err
	}

	return nil
}
