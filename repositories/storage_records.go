package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type StorageRecords interface {
	Create(data entities.StorageRecord) error
	FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error)
	FindManyByLocation(location entities.Location, offset, limit int64) (int64, []entities.StorageRecord, error)
	UpdateOneByID(id string, data entities.StorageRecord) error
}

type mongoStorageRecords struct {
	collection *mongo.Collection
}

func NewMongoStorageRecords(database *mongo.Database) StorageRecords {
	return &mongoStorageRecords{
		collection: database.Collection("storage_records"),
	}
}

func (repository *mongoStorageRecords) Create(data entities.StorageRecord) error {
	model := models.NewStorageRecord(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return err
	}

	return nil
}

func (repository *mongoStorageRecords) FindOneByProductTypeIDAndLocation(productTypeID string, location entities.Location) (entities.StorageRecord, error) {
	var model models.StorageRecord

	filter := bson.M{
		"$and": bson.A{
			bson.M{"product_type_id": productTypeID},
			bson.M{"location.id": location.ID},
			bson.M{"location.type": location.Type},
		},
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.StorageRecord{}, nil
		}

		return entities.StorageRecord{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoStorageRecords) FindManyByLocation(location entities.Location, offset, limit int64) (int64, []entities.StorageRecord, error) {
	modelList := make([]models.FindByLocationStorageRecord, 0)
	entityList := make([]entities.StorageRecord, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{"location.id": location.ID},
			bson.M{"location.type": location.Type},
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
			{"$sort", bson.M{"quantity": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "product_types"},
				{"localField", "product_type_id"},
				{"foreignField", "_id"},
				{"as", "product_type"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$product_type"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoStorageRecords) UpdateOneByID(id string, data entities.StorageRecord) error {
	model := models.NewUpdatedStorageRecord(data)

	update := bson.M{
		"$set": &model,
	}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
