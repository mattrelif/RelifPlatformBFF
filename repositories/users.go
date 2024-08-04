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

type Users interface {
	CreateUser(data entities.User) (entities.User, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	CountByEmail(email string) (int64, error)
	CountById(email string) (int64, error)
	UpdateOneById(id string, data entities.User) error
}

type usersMongo struct {
	collection *mongo.Collection
}

func NewUsersMongo(database *mongo.Database) Users {
	return &usersMongo{
		collection: database.Collection("users"),
	}
}

func (repository *usersMongo) CreateUser(data entities.User) (entities.User, error) {
	model := models.NewUser(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error) {
	modelList := make([]models.User, 0)
	entityList := make([]entities.User, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"organization_id": organizationId,
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

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"first_name": 1})
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

func (repository *usersMongo) FindOneById(id string) (entities.User, error) {
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

func (repository *usersMongo) FindOneByEmail(email string) (entities.User, error) {
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

func (repository *usersMongo) CountByEmail(email string) (int64, error) {
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

func (repository *usersMongo) CountById(id string) (int64, error) {
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

func (repository *usersMongo) UpdateOneById(id string, data entities.User) error {
	model := models.NewUpdatedUser(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
