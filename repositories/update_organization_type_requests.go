package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type UpdateOrganizationTypeRequests interface {
	Create(data entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error)
	FindOneById(id string) (entities.UpdateOrganizationTypeRequest, error)
	FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	UpdateOneById(id string, data entities.UpdateOrganizationTypeRequest) error
	CountPendingByOrganizationId(organizationId string) (int64, error)
}

type mongoUpdateOrganizationTypeRequests struct {
	collection *mongo.Collection
}

func NewMongoUpdateOrganizationTypeRequests(database *mongo.Database) UpdateOrganizationTypeRequests {
	return &mongoUpdateOrganizationTypeRequests{
		collection: database.Collection("update_organization_type_requests"),
	}
}

func (repository *mongoUpdateOrganizationTypeRequests) Create(data entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error) {
	model := models.NewUpdateOrganizationTypeRequest(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUpdateOrganizationTypeRequests) FindOneById(id string) (entities.UpdateOrganizationTypeRequest, error) {
	var model models.UpdateOrganizationTypeRequest

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.UpdateOrganizationTypeRequest{}, utils.ErrUpdateOrganizationTypeRequestNotFound
		}
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUpdateOrganizationTypeRequests) FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	modelList := make([]models.FindUpdateOrganizationTypeRequest, 0)
	entityList := make([]entities.UpdateOrganizationTypeRequest, 0)

	filter := bson.M{}

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
			{"lookup", bson.D{
				{"from", "organizations"},
				{"localField", "organization_id"},
				{"foreignKey", "_id"},
				{"as", "organization"},
			}},
		},
		bson.D{
			{"lookup", bson.D{
				{"from", "users"},
				{"localField", "creator_id"},
				{"foreignKey", "_id"},
				{"as", "creator"},
			}},
		},
		bson.D{
			{"lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignKey", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$creator"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoUpdateOrganizationTypeRequests) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	modelList := make([]models.FindUpdateOrganizationTypeRequest, 0)
	entityList := make([]entities.UpdateOrganizationTypeRequest, 0)

	filter := bson.M{"organization_id": organizationId}

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
			{"lookup", bson.D{
				{"from", "users"},
				{"localField", "creator_id"},
				{"foreignKey", "_id"},
				{"as", "creator"},
			}},
		},
		bson.D{
			{"lookup", bson.D{
				{"from", "users"},
				{"localField", "auditor_id"},
				{"foreignKey", "_id"},
				{"as", "auditor"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$creator"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoUpdateOrganizationTypeRequests) UpdateOneById(id string, data entities.UpdateOrganizationTypeRequest) error {
	model := models.NewUpdatedUpdateOrganizationTypeRequest(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoUpdateOrganizationTypeRequests) CountPendingByOrganizationId(organizationId string) (int64, error) {
	filter := bson.M{"organization_id": organizationId, "status": utils.PendingStatus}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}
