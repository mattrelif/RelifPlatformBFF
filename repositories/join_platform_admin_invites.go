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

type JoinPlatformAdminInvites interface {
	Create(invite entities.JoinPlatformAdminInvite) (entities.JoinPlatformAdminInvite, error)
	FindManyPaginated(offset, limit int64) (int64, []entities.JoinPlatformAdminInvite, error)
	FindOneAndDeleteByCode(code string) (entities.JoinPlatformAdminInvite, error)
}

type joinPlatformAdminInvites struct {
	collection *mongo.Collection
}

func NewJoinPlatformAdminInvites(database *mongo.Database) JoinPlatformAdminInvites {
	return &joinPlatformAdminInvites{
		collection: database.Collection("join_platform_admin_invites"),
	}
}

func (repository *joinPlatformAdminInvites) Create(invite entities.JoinPlatformAdminInvite) (entities.JoinPlatformAdminInvite, error) {
	model := models.NewJoinPlatformAdminInvite(invite)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.JoinPlatformAdminInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *joinPlatformAdminInvites) FindManyPaginated(offset, limit int64) (int64, []entities.JoinPlatformAdminInvite, error) {
	entityList := make([]entities.JoinPlatformAdminInvite, 0)

	count, err := repository.collection.CountDocuments(context.Background(), bson.M{})

	if err != nil {
		return 0, nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$sort", bson.M{"created_at": 1}},
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
				{"localField", "inviter_id"},
				{"foreignField", "_id"},
				{"as", "inviter"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$inviter"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.FindJoinPlatformAdminInvite

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *joinPlatformAdminInvites) FindOneAndDeleteByCode(code string) (entities.JoinPlatformAdminInvite, error) {
	var model models.JoinPlatformAdminInvite

	filter := bson.M{"code": code}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.JoinPlatformAdminInvite{}, utils.ErrJoinPlatformAdminInviteNotFound
		}
		return entities.JoinPlatformAdminInvite{}, err
	}

	return model.ToEntity(), nil
}
