package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
)

type ProductTypeAllocations interface {
	Create(data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error)
	FindManyByProductTypeIDPaginated(productTypeID string, offset, limit int64) (int64, []entities.ProductTypeAllocation, error)
	DeleteManyByProductTypeID(productTypeID string) error
}

type mongoProductTypeAllocations struct {
	collection *mongo.Collection
}

func NewMongoProductTypeAllocations(database *mongo.Database) ProductTypeAllocations {
	return &mongoProductTypeAllocations{
		collection: database.Collection("product_type_allocations"),
	}
}

func (repository *mongoProductTypeAllocations) Create(data entities.ProductTypeAllocation) (entities.ProductTypeAllocation, error) {
	model := models.NewProductTypeAllocation(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.ProductTypeAllocation{}, err
	}

	return data, nil
}

func (repository *mongoProductTypeAllocations) FindManyByProductTypeIDPaginated(productTypeID string, offset, limit int64) (int64, []entities.ProductTypeAllocation, error) {
	entityList := make([]entities.ProductTypeAllocation, 0)

	filter := bson.M{"product_type_id": productTypeID}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{{"$match", filter}},
		bson.D{{"$sort", bson.M{"created_at": -1}}},
		bson.D{{"$skip", offset}},
		bson.D{{"$limit", limit}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "location.id"},
				{"foreignField", "_id"},
				{"as", "organization"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$organization"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "product_types"},
				{"localField", "product_type_id"},
				{"foreignField", "_id"},
				{"as", "product_type"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$product_type"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "from.id"},
				{"foreignField", "_id"},
				{"as", "from_housing"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$from_housing"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "to.id"},
				{"foreignField", "_id"},
				{"as", "to_housing"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$to_housing"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$project", bson.M{
				"_id":             1,
				"product_type_id": 1,
				"product_type":    1,
				"type":            1,
				"from": bson.D{
					{"id", "$from.id"},
					{"type", "$from.type"},
					{"name", bson.D{
						{"$cond", bson.D{
							{"if", bson.D{{"$eq", bson.A{"$from.type", utils.OrganizationLocationType}}}},
							{"then", "$organization.name"},
							{"else", "$from_housing.name"},
						}},
					}},
				},
				"to": bson.D{
					{"id", "$to.id"},
					{"type", "$to.type"},
					{"name", bson.D{
						{"$cond", bson.D{
							{"if", bson.D{{"$eq", bson.A{"$to.type", utils.OrganizationLocationType}}}},
							{"then", "$organization.name"},
							{"else", "$to_housing.name"},
						}},
					}},
				},
				"organization_id": 1,
				"quantity":        1,
				"created_at":      1,
			},
		}},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.FindProductTypeAllocation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoProductTypeAllocations) DeleteManyByProductTypeID(productTypeID string) error {
	filter := bson.M{"product_type_id": productTypeID}

	if _, err := repository.collection.DeleteMany(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
