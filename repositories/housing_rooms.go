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

type HousingRooms interface {
	CreateMany(data []entities.HousingRoom) ([]entities.HousingRoom, error)
	FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error)
	FindOneById(id string) (entities.HousingRoom, error)
	FindOneAndUpdateById(id string, data entities.HousingRoom) (entities.HousingRoom, error)
	IncreaseAvailableVacanciesById(id string) error
	DecreaseAvailableVacanciesById(id string) error
	DeleteOneById(id string) error
}

type mongoHousingRooms struct {
	collection *mongo.Collection
}

func NewMongoHousingRooms(database *mongo.Database) HousingRooms {
	return &mongoHousingRooms{
		collection: database.Collection("housing_rooms"),
	}
}

func (repository *mongoHousingRooms) CreateMany(data []entities.HousingRoom) ([]entities.HousingRoom, error) {
	modelList := make([]interface{}, 0)
	entityList := make([]entities.HousingRoom, 0)

	for _, room := range data {
		model := models.HousingRoom{
			ID:                 primitive.NewObjectID().Hex(),
			Name:               room.Name,
			HousingID:          room.HousingID,
			Status:             room.Status,
			TotalVacancies:     room.TotalVacancies,
			AvailableVacancies: room.TotalVacancies,
			CreatedAt:          time.Now(),
		}

		modelList = append(modelList, model)
		entityList = append(entityList, model.ToEntity())
	}

	if _, err := repository.collection.InsertMany(context.Background(), modelList); err != nil {
		return nil, err
	}

	return entityList, nil
}

func (repository *mongoHousingRooms) FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.HousingRoom, error) {
	modelList := make([]models.HousingRoom, 0)
	entityList := make([]entities.HousingRoom, 0)

	filter := bson.M{"housing_id": housingId}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"name": 1})
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

func (repository *mongoHousingRooms) FindOneById(id string) (entities.HousingRoom, error) {
	var model models.HousingRoom

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.HousingRoom{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousingRooms) FindOneAndUpdateById(id string, data entities.HousingRoom) (entities.HousingRoom, error) {
	model := models.HousingRoom{
		Name:           data.Name,
		TotalVacancies: data.TotalVacancies,
		UpdatedAt:      time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.HousingRoom{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousingRooms) IncreaseAvailableVacanciesById(id string) error {
	model := models.HousingRoom{
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model, "$inc": bson.M{"available_vacancies": 1}}

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *mongoHousingRooms) DecreaseAvailableVacanciesById(id string) error {
	model := models.HousingRoom{
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model, "$inc": bson.M{"available_vacancies": -1}}

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *mongoHousingRooms) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
