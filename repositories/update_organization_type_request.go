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

type UpdateOrganizationTypeRequests interface {
	Create(request entities.UpdateOrganizationTypeRequest) (string, error)
	FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindOneAndUpdateById(id string, request entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error)
}

type mongoUpdateOrganizationTypeRequests struct {
	collection *mongo.Collection
}

func NewMongoUpdateOrganizationTypeRequests(database *mongo.Database) UpdateOrganizationTypeRequests {
	return &mongoUpdateOrganizationTypeRequests{
		collection: database.Collection("update_organization_type_requests"),
	}
}

func (repository *mongoUpdateOrganizationTypeRequests) Create(request entities.UpdateOrganizationTypeRequest) (string, error) {
	model := models.UpdateOrganizationTypeRequest{
		ID:             primitive.NewObjectID().Hex(),
		OrganizationID: request.OrganizationID,
		CreatorID:      request.CreatorID,
		Status:         request.Status,
		CreatedAt:      request.CreatedAt,
	}

	result, err := repository.collection.InsertOne(context.Background(), &model)

	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repository *mongoUpdateOrganizationTypeRequests) FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	modelList := make([]models.UpdateOrganizationTypeRequest, 0)
	entityList := make([]entities.UpdateOrganizationTypeRequest, 0)

	count, err := repository.collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)
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

func (repository *mongoUpdateOrganizationTypeRequests) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	modelList := make([]models.UpdateOrganizationTypeRequest, 0)
	entityList := make([]entities.UpdateOrganizationTypeRequest, 0)

	oid, err := primitive.ObjectIDFromHex(organizationId)

	if err != nil {
		return 0, nil, err
	}

	filter := bson.M{"organization_id": oid}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)

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

func (repository *mongoUpdateOrganizationTypeRequests) FindOneAndUpdateById(id string, request entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error) {
	model := models.UpdateOrganizationTypeRequest{
		AuditorID:    request.AuditorID,
		Status:       request.Status,
		RejectReason: request.RejectReason,
		RejectedAt:   request.RejectedAt,
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	return model.ToEntity(), nil
}
