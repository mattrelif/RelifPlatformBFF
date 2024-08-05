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

type JoinOrganizationRequests interface {
	Create(data entities.JoinOrganizationRequest) (entities.JoinOrganizationRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneById(id string) (entities.JoinOrganizationRequest, error)
	UpdateOneById(id string, data entities.JoinOrganizationRequest) error
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

func (repository *mongoJoinOrganizationRequests) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	modelList := make([]models.JoinOrganizationRequest, 0)
	entityList := make([]entities.JoinOrganizationRequest, 0)

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

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoJoinOrganizationRequests) FindOneById(id string) (entities.JoinOrganizationRequest, error) {
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

func (repository *mongoJoinOrganizationRequests) UpdateOneById(id string, data entities.JoinOrganizationRequest) error {
	model := models.NewUpdatedJoinOrganizationRequest(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
