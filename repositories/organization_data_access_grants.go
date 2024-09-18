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

type OrganizationDataAccessGrants interface {
	Create(data entities.OrganizationDataAccessGrant) error
	FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindManyByTargetOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error)
	FindOneByID(id string) (entities.OrganizationDataAccessGrant, error)
	DeleteOneByID(id string) error
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

func (repository *mongoOrganizationDataAccessGrants) FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	entityList := make([]entities.OrganizationDataAccessGrant, 0)

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
			{"$sort", bson.M{"created_at": -1}},
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
				{"localField", "target_organization_id"},
				{"foreignField", "_id"},
				{"as", "target_organization"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$target_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}
	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.OrganizationDataAccessGrant

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizationDataAccessGrants) FindManyByTargetOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessGrant, error) {
	entityList := make([]entities.OrganizationDataAccessGrant, 0)

	filter := bson.M{"target_organization_id": organizationID}

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

	for cursor.Next(context.Background()) {
		var model models.OrganizationDataAccessGrant

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoOrganizationDataAccessGrants) FindOneByID(id string) (entities.OrganizationDataAccessGrant, error) {
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

func (repository *mongoOrganizationDataAccessGrants) DeleteOneByID(id string) error {
	filter := bson.M{"_id": id}

	if _, err := repository.collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
