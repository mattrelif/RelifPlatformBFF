package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
	"time"
)

type HousingRooms interface {
	CreateMany(data []entities.HousingRoom, housingID string) ([]entities.HousingRoom, error)
	FindManyByHousingID(housingID string, limit, offset int64) (int64, []entities.HousingRoom, error)
	FindOneByID(id string) (entities.HousingRoom, error)
	FindOneCompleteByID(id string) (entities.HousingRoom, error)
	UpdateOneByID(id string, data entities.HousingRoom) error
	IncreaseAvailableVacanciesByID(id string) error
	DecreaseAvailableVacanciesByID(id string) error
}

type mongoHousingRooms struct {
	collection *mongo.Collection
}

func NewMongoHousingRooms(database *mongo.Database) HousingRooms {
	return &mongoHousingRooms{
		collection: database.Collection("housing_rooms"),
	}
}

func (repository *mongoHousingRooms) CreateMany(data []entities.HousingRoom, housingID string) ([]entities.HousingRoom, error) {
	modelList := make([]interface{}, 0)
	entityList := make([]entities.HousingRoom, 0)

	for _, room := range data {
		room.HousingID = housingID
		model := models.NewHousingRoom(room)

		modelList = append(modelList, model)
		entityList = append(entityList, model.ToEntity())
	}

	if _, err := repository.collection.InsertMany(context.Background(), modelList); err != nil {
		return nil, err
	}

	return entityList, nil
}

func (repository *mongoHousingRooms) FindManyByHousingID(housingID string, limit, offset int64) (int64, []entities.HousingRoom, error) {
	modelList := make([]models.FindHousingRoom, 0)
	entityList := make([]entities.HousingRoom, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"housing_id": housingID,
			},
			bson.M{
				"status": bson.M{
					"$ne": utils.InactiveStatus,
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
				{"from", "beneficiaries"},
				{"let", bson.D{{"roomID", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$current_room_id", "$$roomID"}}},
									bson.D{{"$ne", bson.A{"status", utils.InactiveStatus}}},
								}},
							}},
						}},
					},
				}},
				{"as", "beneficiaries"},
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

func (repository *mongoHousingRooms) FindOneByID(id string) (entities.HousingRoom, error) {
	var model models.HousingRoom

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$ne": utils.InactiveStatus,
				},
			},
		},
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.HousingRoom{}, utils.ErrHousingRoomNotFound
		}

		return entities.HousingRoom{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousingRooms) FindOneCompleteByID(id string) (entities.HousingRoom, error) {
	var model models.FindHousingRoom

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"_id": id,
			},
			bson.M{
				"status": bson.M{
					"$ne": utils.InactiveStatus,
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
				{"from", "beneficiaries"},
				{"let", bson.D{{"roomID", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$expr", bson.D{
								{"$and", bson.A{
									bson.D{{"$eq", bson.A{"$current_room_id", "$$roomID"}}},
									bson.D{{"$ne", bson.A{"status", utils.InactiveStatus}}},
								}},
							}},
						}},
					},
				}},
				{"as", "beneficiaries"},
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
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return entities.HousingRoom{}, err
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.HousingRoom{}, err
		}
	} else {
		return entities.HousingRoom{}, utils.ErrHousingRoomNotFound
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousingRooms) UpdateOneByID(id string, data entities.HousingRoom) error {
	model := models.NewUpdatedHousingRoom(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoHousingRooms) IncreaseAvailableVacanciesByID(id string) error {
	model := models.HousingRoom{
		UpdatedAt: time.Now(),
	}

	update := bson.M{"$set": &model, "$inc": bson.M{"available_vacancies": 1}}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoHousingRooms) DecreaseAvailableVacanciesByID(id string) error {
	model := models.HousingRoom{
		UpdatedAt: time.Now(),
	}

	update := bson.M{"$set": &model, "$inc": bson.M{"available_vacancies": -1}}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
