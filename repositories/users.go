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

type Users interface {
	CreateUser(data entities.User) (entities.User, error)
	FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.User, error)
	FindOneByID(id string) (entities.User, error)
	FindOneCompleteByID(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	CountByEmail(email string) (int64, error)
	CountByID(email string) (int64, error)
	UpdateOneByID(id string, data entities.User) error
}

type mongoUsers struct {
	collection *mongo.Collection
}

func NewUsersMongo(database *mongo.Database) Users {
	return &mongoUsers{
		collection: database.Collection("users"),
	}
}

func (repository *mongoUsers) CreateUser(data entities.User) (entities.User, error) {
	model := models.NewUser(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUsers) FindManyByOrganizationID(organizationID string, offset, limit int64) (int64, []entities.User, error) {
	modelList := make([]models.FindUser, 0)
	entityList := make([]entities.User, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"organization_id": organizationID,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$sort", bson.M{"first_name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
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
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoUsers) FindOneByID(id string) (entities.User, error) {
	var model models.User

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, utils.ErrUserNotFound
		}

		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUsers) FindOneCompleteByID(id string) (entities.User, error) {
	var model models.FindUser

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
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
			{"$unwind", bson.D{{"path", "$organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return entities.User{}, err
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.User{}, err
		}
	} else {
		return entities.User{}, utils.ErrUserNotFound
	}

	return model.ToEntity(), nil
}

func (repository *mongoUsers) FindOneByEmail(email string) (entities.User, error) {
	var model models.User

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"email": email,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.User{}, utils.ErrUserNotFound
		}

		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoUsers) CountByEmail(email string) (int64, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"email": email,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *mongoUsers) CountByID(id string) (int64, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$not": bson.M{
						"$eq": utils.InactiveStatus,
					},
				},
			},
		},
	}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *mongoUsers) UpdateOneByID(id string, data entities.User) error {
	model := models.NewUpdatedUser(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
