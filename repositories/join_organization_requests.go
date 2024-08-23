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

type JoinOrganizationRequests interface {
	Create(data entities.JoinOrganizationRequest) (entities.JoinOrganizationRequest, error)
	FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindManyByUserIDPaginated(userID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneByID(id string) (entities.JoinOrganizationRequest, error)
	UpdateOneByID(id string, data entities.JoinOrganizationRequest) error
}

type mongoJoinOrganizationRequests struct {
	collection *mongo.Collection
}

func NewMongoJoinOrganizationRequests(database *mongo.Database) JoinOrganizationRequests {
	return &mongoJoinOrganizationRequests{
		collection: database.Collection("join_organization_requests"),
	}
}

func (repository *mongoJoinOrganizationRequests) Create(data entities.JoinOrganizationRequest) (entities.JoinOrganizationRequest, error) {
	model := models.NewJoinOrganizationRequest(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.JoinOrganizationRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationRequests) FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	modelList := make([]models.FindJoinOrganizationRequest, 0)
	entityList := make([]entities.JoinOrganizationRequest, 0)

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
				{"from", "users"},
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "user"},
			}},
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
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignField", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoJoinOrganizationRequests) FindManyByUserIDPaginated(userID string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	modelList := make([]models.FindJoinOrganizationRequest, 0)
	entityList := make([]entities.JoinOrganizationRequest, 0)

	filter := bson.M{"user_id": userID}
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
				{"localField", "user_id"},
				{"foreignField", "_id"},
				{"as", "user"},
			}},
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
			{"$lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignField", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoJoinOrganizationRequests) FindOneByID(id string) (entities.JoinOrganizationRequest, error) {
	var model models.JoinOrganizationRequest

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.JoinOrganizationRequest{}, utils.ErrJoinOrganizationRequestNotFound
		}
		return entities.JoinOrganizationRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationRequests) UpdateOneByID(id string, data entities.JoinOrganizationRequest) error {
	model := models.NewUpdatedJoinOrganizationRequest(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
