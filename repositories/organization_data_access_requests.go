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

type OrganizationDataAccessRequests interface {
	Create(data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByTargetOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindOneByID(id string) (entities.OrganizationDataAccessRequest, error)
	UpdateOneByID(id string, data entities.OrganizationDataAccessRequest) error
}

type mongoOrganizationDataAccessRequests struct {
	collection *mongo.Collection
}

func NewMongoOrganizationDataAccessRequests(database *mongo.Database) OrganizationDataAccessRequests {
	return &mongoOrganizationDataAccessRequests{
		collection: database.Collection("access_organization_data_requests"),
	}
}

func (repository *mongoOrganizationDataAccessRequests) Create(data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error) {
	model := models.NewOrganizationDataAccessRequest(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizationDataAccessRequests) FindManyByRequesterOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	modelList := make([]models.FindOrganizationDataAccessRequest, 0)
	entityList := make([]entities.OrganizationDataAccessRequest, 0)

	filter := bson.M{"requester_organization_id": organizationID}
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
				{"from", "users"},
				{"localField", "requester_id"},
				{"foreignField", "_id"},
				{"as", "requester"},
			}},
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
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignField", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$requester"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$target_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$auditor"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoOrganizationDataAccessRequests) FindManyByTargetOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	modelList := make([]models.FindOrganizationDataAccessRequest, 0)
	entityList := make([]entities.OrganizationDataAccessRequest, 0)

	filter := bson.M{"target_organization_id": organizationID}
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
				{"from", "users"},
				{"localField", "requester_id"},
				{"foreignField", "_id"},
				{"as", "requester"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "requester_organization_id"},
				{"foreignField", "_id"},
				{"as", "requester_organization"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignField", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$requester"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$requester_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$auditor"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoOrganizationDataAccessRequests) UpdateOneByID(id string, data entities.OrganizationDataAccessRequest) error {
	model := models.NewUpdatedOrganizationDataAccessRequest(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoOrganizationDataAccessRequests) FindOneByID(id string) (entities.OrganizationDataAccessRequest, error) {
	var model models.OrganizationDataAccessRequest

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.OrganizationDataAccessRequest{}, utils.ErrOrganizationDataAccessRequestNotFound
		}
		return entities.OrganizationDataAccessRequest{}, err
	}

	return model.ToEntity(), nil
}
