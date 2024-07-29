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

type JoinOrganizationInvites interface {
	Create(data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	FindOneAndDeleteById(id string) (entities.JoinOrganizationInvite, error)
	DeleteOneById(id string) error
}

type mongoJoinOrganizationInvites struct {
	collection *mongo.Collection
}

func NewMongoJoinOrganizationInvites(database *mongo.Database) JoinOrganizationInvites {
	return &mongoJoinOrganizationInvites{
		collection: database.Collection("join_organization_invites"),
	}
}

func (repository *mongoJoinOrganizationInvites) Create(invite entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error) {
	model := models.JoinOrganizationInvite{
		ID:             primitive.NewObjectID().Hex(),
		UserID:         invite.UserID,
		OrganizationID: invite.OrganizationID,
		CreatorID:      invite.CreatorID,
		CreatedAt:      time.Now(),
		ExpiresAt:      time.Now().Add(4 * time.Hour),
	}

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.JoinOrganizationInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationInvites) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	modelList := make([]models.JoinOrganizationInvite, 0)
	entityList := make([]entities.JoinOrganizationInvite, 0)

	filter := bson.M{"organization_id": organizationId}
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

func (repository *mongoJoinOrganizationInvites) FindOneAndDeleteById(id string) (entities.JoinOrganizationInvite, error) {
	var model models.JoinOrganizationInvite

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Decode(&model); err != nil {
		return entities.JoinOrganizationInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationInvites) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
