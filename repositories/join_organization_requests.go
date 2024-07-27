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

type JoinOrganizationRequests interface {
	Create(request entities.JoinOrganizationRequest) (string, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error)
	FindOneAndDeleteById(id string) (entities.JoinOrganizationRequest, error)
	DeleteOneById(id string) error
}

type mongoJoinOrganizationRequests struct {
	collection *mongo.Collection
}

func NewMongoJoinOrganizationRequests(database *mongo.Database) JoinOrganizationRequests {
	return &mongoJoinOrganizationRequests{
		collection: database.Collection("join_organization_requests"),
	}
}

func (repository *mongoJoinOrganizationRequests) Create(request entities.JoinOrganizationRequest) (string, error) {
	model := models.JoinOrganizationRequest{
		ID:             primitive.NewObjectID().Hex(),
		UserID:         request.UserID,
		OrganizationID: request.OrganizationID,
		CreatedAt:      request.CreatedAt,
		ExpiresAt:      request.ExpiresAt,
	}

	result, err := repository.collection.InsertOne(context.Background(), model)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repository *mongoJoinOrganizationRequests) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationRequest, error) {
	modelList := make([]models.JoinOrganizationRequest, 0)
	entityList := make([]entities.JoinOrganizationRequest, 0)

	oid, err := primitive.ObjectIDFromHex(organizationId)

	if err != nil {
		return 0, nil, err
	}

	filter := bson.M{"organization_id": oid}

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

func (repository *mongoJoinOrganizationRequests) FindOneAndDeleteById(id string) (entities.JoinOrganizationRequest, error) {
	var model models.JoinOrganizationRequest

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.JoinOrganizationRequest{}, err
	}

	filter := bson.M{"_id": oid}

	if err = repository.collection.FindOneAndDelete(context.Background(), filter).Decode(&model); err != nil {
		return entities.JoinOrganizationRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationRequests) DeleteOneById(id string) error {
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
