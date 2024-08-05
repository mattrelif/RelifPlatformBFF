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

type JoinOrganizationInvites interface {
	Create(data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error)
	FindOneById(id string) (entities.JoinOrganizationInvite, error)
	UpdateOneById(id string, data entities.JoinOrganizationInvite) error
}

type mongoJoinOrganizationInvites struct {
	collection *mongo.Collection
}

func NewMongoJoinOrganizationInvites(database *mongo.Database) JoinOrganizationInvites {
	return &mongoJoinOrganizationInvites{
		collection: database.Collection("join_organization_invites"),
	}
}

func (repository *mongoJoinOrganizationInvites) Create(data entities.JoinOrganizationInvite) (entities.JoinOrganizationInvite, error) {
	model := models.NewJoinOrganizationInvite(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.JoinOrganizationInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationInvites) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	modelList := make([]models.FindJoinOrganizationInvite, 0)
	entityList := make([]entities.JoinOrganizationInvite, 0)

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
				{"localField", "creator_id"},
				{"foreignField", "_id"},
				{"as", "creator"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$creator"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoJoinOrganizationInvites) FindManyByUserId(userId string, offset, limit int64) (int64, []entities.JoinOrganizationInvite, error) {
	modelList := make([]models.FindJoinOrganizationInvite, 0)
	entityList := make([]entities.JoinOrganizationInvite, 0)

	filter := bson.M{"user_id": userId}
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
				{"localField", "creator_id"},
				{"foreignField", "_id"},
				{"as", "creator"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$user"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$creator"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoJoinOrganizationInvites) FindOneById(id string) (entities.JoinOrganizationInvite, error) {
	var model models.JoinOrganizationInvite

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.JoinOrganizationInvite{}, utils.ErrJoinOrganizationInviteNotFound
		}

		return entities.JoinOrganizationInvite{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoJoinOrganizationInvites) UpdateOneById(id string, data entities.JoinOrganizationInvite) error {
	model := models.NewUpdatedJoinOrganizationInvite(data)

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateOne(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}
