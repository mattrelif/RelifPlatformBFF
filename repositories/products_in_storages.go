package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type ProductsInStorages interface {
	CreateMany(data entities.ProductInStorage, quantity int) error
	FindManyIDsByLocation(location entities.Location, quantity int) ([]interface{}, error)
	UpdateManyByIDs(ids []interface{}, data entities.ProductInStorage) error
	DeleteManyByIDs(ids []interface{}) error
}

type mongoProductsInStorages struct {
	collection *mongo.Collection
}

func NewMongoProductsInStorages(database *mongo.Database) ProductsInStorages {
	return &mongoProductsInStorages{
		collection: database.Collection("products_in_storages"),
	}
}

func (repository *mongoProductsInStorages) CreateMany(data entities.ProductInStorage, quantity int) error {
	modelList := make([]interface{}, quantity)

	for i := 0; i < quantity; i++ {
		modelList[i] = models.NewProductInStorage(data)
	}

	if _, err := repository.collection.InsertMany(context.Background(), modelList); err != nil {
		return err
	}

	return nil
}

func (repository *mongoProductsInStorages) FindManyIDsByLocation(location entities.Location, quantity int) ([]interface{}, error) {
	ids := make([]interface{}, 0)

	filter := bson.M{"location.id": location.ID, "location.type": location.Type}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$project", bson.M{"_id": 1}},
		},
		bson.D{
			{"$sample", bson.M{"size": quantity}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result bson.M

		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}

		ids = append(ids, result["_id"])
	}

	return ids, nil
}

func (repository *mongoProductsInStorages) UpdateManyByIDs(ids []interface{}, data entities.ProductInStorage) error {
	model := models.NewUpdatedProductInStorage(data)

	filter := bson.M{"_id": bson.M{"$in": ids}}
	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateMany(context.Background(), filter, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoProductsInStorages) DeleteManyByIDs(ids []interface{}) error {
	filter := bson.M{"_id": bson.M{"$in": ids}}

	if _, err := repository.collection.DeleteMany(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
