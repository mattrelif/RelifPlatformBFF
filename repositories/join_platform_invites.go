package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type JoinPlatformInvites interface {
	Create(data entities.JoinPlatformInvite) (entities.JoinPlatformInvite, error)
	FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.JoinPlatformInvite, error)
	FindOneAndDeleteByCode(code string) (entities.JoinPlatformInvite, error)
}

type mongoJoinPlatformInvites struct {
	collection *mongo.Collection
}

func NewMongoJoinPlatformInvites(database *mongo.Database) JoinPlatformInvites {
	return &mongoJoinPlatformInvites{
		collection: database.Collection("join_platform_invites"),
	}
}

func (repository *mongoJoinPlatformInvites) Create(data entities.JoinPlatformInvite) (entities.JoinPlatformInvite, error) {
	model := models.NewJoinPlatformInvite(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinPlatformInvites) FindManyByOrganizationIDPaginated(organizationID string, offset, limit int64) (int64, []entities.JoinPlatformInvite, error) {
	entityList := make([]entities.JoinPlatformInvite, 0)

	filter := bson.M{"organization_id": organizationID}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.JoinPlatformInvite

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoJoinPlatformInvites) FindOneAndDeleteByCode(code string) (entities.JoinPlatformInvite, error) {
	var model models.JoinPlatformInvite

	filter := bson.M{"code": code}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Decode(&model); err != nil {
		return entities.JoinPlatformInvite{}, err
	}

	return model.ToEntity(), nil
}
