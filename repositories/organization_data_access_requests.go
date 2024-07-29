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

type OrganizationDataAccessRequests interface {
	Create(data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error)
	FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error)
	FindOneAndUpdateById(id string, data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error)
	UpdateOneById(id string, data entities.OrganizationDataAccessRequest) error
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
	model := models.OrganizationDataAccessRequest{
		ID:                      primitive.NewObjectID().Hex(),
		RequesterID:             data.RequesterID,
		RequesterOrganizationID: data.RequesterOrganizationID,
		TargetOrganizationID:    data.TargetOrganizationID,
		Status:                  data.Status,
		CreatedAt:               time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoOrganizationDataAccessRequests) FindManyByRequesterOrganizationId(organizationId string, limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	modelList := make([]models.OrganizationDataAccessRequest, 0)
	entityList := make([]entities.OrganizationDataAccessRequest, 0)

	filter := bson.M{"requester_organization_id": organizationId}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
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

func (repository *mongoOrganizationDataAccessRequests) FindMany(limit, offset int64) (int64, []entities.OrganizationDataAccessRequest, error) {
	modelList := make([]models.OrganizationDataAccessRequest, 0)
	entityList := make([]entities.OrganizationDataAccessRequest, 0)

	count, err := repository.collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
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

func (repository *mongoOrganizationDataAccessRequests) UpdateOneById(id string, data entities.OrganizationDataAccessRequest) error {
	model := models.OrganizationDataAccessRequest{
		AuditorID:    data.AuditorID,
		Status:       data.Status,
		AcceptedAt:   data.AcceptedAt,
		RejectedAt:   data.RejectedAt,
		RejectReason: data.RejectReason,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *mongoOrganizationDataAccessRequests) FindOneAndUpdateById(id string, data entities.OrganizationDataAccessRequest) (entities.OrganizationDataAccessRequest, error) {
	model := models.OrganizationDataAccessRequest{
		AuditorID:    data.AuditorID,
		Status:       data.Status,
		AcceptedAt:   data.AcceptedAt,
		RejectedAt:   data.RejectedAt,
		RejectReason: data.RejectReason,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.OrganizationDataAccessRequest{}, err
	}

	return model.ToEntity(), nil
}
