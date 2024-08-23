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

type Housings interface {
	Create(data entities.Housing) (entities.Housing, error)
	FindManyByOrganizationIDPaginated(organizationID, search string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	FindOneCompleteByID(id string) (entities.Housing, error)
	UpdateOneByID(id string, data entities.Housing) error
	DeleteOneByID(id string) error
}

type mongoHousings struct {
	collection *mongo.Collection
}

func NewMongoHousings(database *mongo.Database) Housings {
	return &mongoHousings{
		collection: database.Collection("housings"),
	}
}

func (repository *mongoHousings) Create(data entities.Housing) (entities.Housing, error) {
	model := models.NewHousing(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindManyByOrganizationIDPaginated(organizationID, search string, limit, offset int64) (int64, []entities.Housing, error) {
	var filter bson.M

	modelList := make([]models.FindHousing, 0)
	entityList := make([]entities.Housing, 0)

	if search != "" {
		filter = bson.M{
			"$and": bson.A{
				bson.M{
					"organization_id": organizationID,
				},
				bson.M{
					"name": bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}
	} else {
		filter = bson.M{
			"organization_id": organizationID,
		}
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
			{"$sort", bson.M{"name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "_id"},
				{"foreignField", "housing_id"},
				{"as", "rooms"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "beneficiaries"},
				{"localField", "_id"},
				{"foreignField", "current_housing_id"},
				{"as", "beneficiaries"},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"total_vacancies", bson.D{
					{"$sum", "$rooms.total_vacancies"},
				}},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"occupied_vacancies", bson.D{
					{"$size", "$beneficiaries"},
				}},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"beneficiaries", 0},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"rooms", 0},
			}},
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

func (repository *mongoHousings) FindOneByID(id string) (entities.Housing, error) {
	var model models.Housing

	filter := bson.M{
		"_id": id,
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Housing{}, utils.ErrHousingNotFound
		}

		return entities.Housing{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) FindOneCompleteByID(id string) (entities.Housing, error) {
	var model models.FindHousing

	filter := bson.M{
		"_id": id,
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "_id"},
				{"foreignField", "housing_id"},
				{"as", "rooms"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "beneficiaries"},
				{"localField", "_id"},
				{"foreignField", "current_housing_id"},
				{"as", "beneficiaries"},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"total_vacancies", bson.D{
					{"$sum", "$rooms.total_vacancies"},
				}},
			}},
		},
		bson.D{
			{"$addFields", bson.D{
				{"occupied_vacancies", bson.D{
					{"$size", "$beneficiaries"},
				}},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"beneficiaries", 0},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"rooms", 0},
			}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return entities.Housing{}, err
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.Housing{}, err
		}
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) UpdateOneByID(id string, data entities.Housing) error {
	model := models.NewUpdatedHousing(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoHousings) DeleteOneByID(id string) error {
	filter := bson.M{"_id": id}

	if _, err := repository.collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
