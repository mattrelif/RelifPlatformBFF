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
			"$addFields", bson.D{
				{"from.name", bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{"$eq": bson.M{"from.type": utils.OrganizationLocationType}},
								"then": "$organization.name",
							},
							bson.M{
								"case": bson.M{"$eq": bson.M{"from.type": utils.HousingLocationType}},
								"then": "$from_housing.name",
							},
						},
						"default": "",
					},
				}},
			},
		}},
		bson.D{{
			"$addFields", bson.D{
				{"to.name", bson.M{
					"$switch": bson.M{
						"branches": bson.A{
							bson.M{
								"case": bson.M{"$eq": bson.M{"to.type": utils.OrganizationLocationType}},
								"then": "$organization.name",
							},
							bson.M{
								"case": bson.M{"$eq": bson.M{"to.type": utils.HousingLocationType}},
								"then": "$to_housing.name",
							},
						},
						"default": "",
					},
				}},
			},
		}},
		bson.D{{
			"$project", bson.M{
				"from_housing": 0,
				"to_housing":   0,
				"organization": 0,
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
