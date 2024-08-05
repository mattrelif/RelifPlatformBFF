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

type OrganizationDataAccessGrants interface {
	Create(data entities.OrganizationDataAccessGrant) error
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindOneById(id string) (entities.OrganizationDataAccessGrant, error)
	DeleteOneById(id string) error
	CountByOrganizationIdAndTargetOrganizationId(organizationId, targetOrganizationId string) (int64, error)
}

type mongoOrganizationDataAccessGrants struct {
	collection *mongo.Collection
}

func NewMongoOrganizationDataAccessGrants(database *mongo.Database) OrganizationDataAccessGrants {
	return &mongoOrganizationDataAccessGrants{
		collection: database.Collection("organization_data_access_grants"),
	}
}

func (repository *mongoOrganizationDataAccessGrants) Create(data entities.OrganizationDataAccessGrant) error {
	model := models.NewOrganizationDataAccessGrant(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return err
	}

	return nil
}

func (repository *mongoOrganizationDataAccessGrants) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	entityList := make([]entities.OrganizationDataAccessGrant, 0)
	modelsList := make([]models.OrganizationDataAccessGrant, 0)

	filter := bson.M{"organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &modelsList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelsList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizationDataAccessGrants) FindManyByTargetOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	entityList := make([]entities.OrganizationDataAccessGrant, 0)
	modelsList := make([]models.OrganizationDataAccessGrant, 0)

	filter := bson.M{"target_organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &modelsList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelsList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizationDataAccessGrants) FindOneById(id string) (entities.OrganizationDataAccessGrant, error) {
	var model models.OrganizationDataAccessGrant

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.OrganizationDataAccessGrant{}, utils.ErrOrganizationDataAccessGrantNotFound
		}
		return entities.OrganizationDataAccessGrant{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizationDataAccessGrants) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if _, err := repository.collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}

func (repository *mongoOrganizationDataAccessGrants) CountByOrganizationIdAndTargetOrganizationId(organizationId, targetOrganizationId string) (int64, error) {
	filter := bson.M{"organization_id": organizationId, "target_organization_id": targetOrganizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}
