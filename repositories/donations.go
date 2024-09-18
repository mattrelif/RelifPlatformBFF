package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
)

type Donations interface {
	Create(data entities.Donation) (entities.Donation, error)
	FindManyByBeneficiaryIDPaginated(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error)
	FindManyByProductTypeIDPaginated(productTypeID string, offset, limit int64) (int64, []entities.Donation, error)
	DeleteManyByProductTypeID(productTypeID string) error
}

type mongoDonations struct {
	collection *mongo.Collection
}

func NewDonations(database *mongo.Database) Donations {
	return &mongoDonations{
		collection: database.Collection("donations"),
	}
}

func (repository *mongoDonations) Create(data entities.Donation) (entities.Donation, error) {
	model := models.NewDonation(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.Donation{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoDonations) FindManyByBeneficiaryIDPaginated(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error) {
	entityList := make([]entities.Donation, 0)

	filter := bson.M{"beneficiary_id": beneficiaryID}

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
				{"from", "beneficiaries"},
				{"localField", "beneficiary_id"},
				{"foreignField", "_id"},
				{"as", "beneficiary"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$beneficiary"},
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
				{"from", "housings"},
				{"localField", "location.id"},
				{"foreignField", "_id"},
				{"as", "housing"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$housing"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$project", bson.M{
				"_id":             1,
				"organization_id": 1,
				"beneficiary_id":  1,
				"beneficiary":     1,
				"from": bson.M{
					"id":   1,
					"type": 1,
					"name": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": bson.A{"$from.type", utils.OrganizationLocationType}},
							"then": "$organization.name",
							"else": bson.M{
								"if":   bson.M{"$eq": bson.A{"$from.type", utils.HousingLocationType}},
								"then": "$housing.name",
							},
						},
					},
				},
				"product_type_id": 1,
				"product_type":    1,
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
		var model models.FindDonation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoDonations) FindManyByProductTypeIDPaginated(productTypeID string, offset, limit int64) (int64, []entities.Donation, error) {
	entityList := make([]entities.Donation, 0)

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
				{"from", "beneficiaries"},
				{"localField", "beneficiary_id"},
				{"foreignField", "_id"},
				{"as", "beneficiary"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$beneficiary"},
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
				{"from", "housings"},
				{"localField", "location.id"},
				{"foreignField", "_id"},
				{"as", "housing"},
			},
		}},
		bson.D{{
			"$unwind", bson.D{
				{"path", "$housing"},
				{"preserveNullAndEmptyArrays", true},
			},
		}},
		bson.D{{
			"$project", bson.M{
				"_id":             1,
				"organization_id": 1,
				"beneficiary_id":  1,
				"beneficiary":     1,
				"from": bson.M{
					"id":   1,
					"type": 1,
					"name": bson.M{
						"$cond": bson.M{
							"if":   bson.M{"$eq": bson.A{"$from.type", utils.OrganizationLocationType}},
							"then": "$organization.name",
							"else": bson.M{
								"if":   bson.M{"$eq": bson.A{"$from.type", utils.HousingLocationType}},
								"then": "$housing.name",
							},
						},
					},
				},
				"product_type_id": 1,
				"product_type":    1,
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
		var model models.FindDonation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoDonations) DeleteManyByProductTypeID(productTypeID string) error {
	filter := bson.M{"product_type_id": productTypeID}

	if _, err := repository.collection.DeleteMany(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
