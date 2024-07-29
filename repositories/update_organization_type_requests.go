package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"time"
)

type UpdateOrganizationTypeRequests interface {
	Create(data entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error)
	FindMany(offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error)
	FindOneAndUpdateById(id string, data entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error)
	UpdateOneById(id string, data entities.UpdateOrganizationTypeRequest) error
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
	model := models.UpdateOrganizationTypeRequest{
		ID:             primitive.NewObjectID().Hex(),
		OrganizationID: data.OrganizationID,
		CreatorID:      data.CreatorID,
		Status:         data.Status,
		CreatedAt:      time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	return model.ToEntity(), nil
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

	defer cursor.Close(context.Background())

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

func (repository *mongoUpdateOrganizationTypeRequests) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.UpdateOrganizationTypeRequest, error) {
	modelList := make([]models.UpdateOrganizationTypeRequest, 0)
	entityList := make([]entities.UpdateOrganizationTypeRequest, 0)

	filter := bson.M{"organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit)
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	defer cursor.Close(context.Background())

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

func (repository *mongoUpdateOrganizationTypeRequests) FindOneAndUpdateById(id string, data entities.UpdateOrganizationTypeRequest) (entities.UpdateOrganizationTypeRequest, error) {
	model := models.UpdateOrganizationTypeRequest{
		AuditorID:    data.AuditorID,
		Status:       data.Status,
		RejectReason: data.RejectReason,
		RejectedAt:   data.RejectedAt,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.UpdateOrganizationTypeRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUpdateOrganizationTypeRequests) UpdateOneById(id string, data entities.UpdateOrganizationTypeRequest) error {
	model := models.UpdateOrganizationTypeRequest{
		AuditorID:    data.AuditorID,
		Status:       data.Status,
		RejectReason: data.RejectReason,
		RejectedAt:   data.RejectedAt,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}

	return nil
}
