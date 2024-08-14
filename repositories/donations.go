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
	FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error)
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

func (repository *mongoDonations) FindManyByBeneficiaryID(beneficiaryID string, offset, limit int64) (int64, []entities.Donation, error) {
	modelList := make([]models.FindDonation, 0)
	entityList := make([]entities.Donation, 0)

	filter := bson.M{"beneficiary_id": beneficiaryID}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$sort", bson.M{"created_at": -1}},
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
				{"localField", "beneficiary_id"},
				{"foreignField", "_id"},
				{"as", "beneficiary"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$beneficiary"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$facet", bson.D{
				{"organizationDonations", bson.A{
					bson.D{
						{"$match", bson.M{"location.type": utils.OrganizationLocationType}},
					},
					bson.D{
						{"$lookup", bson.D{
							{"from", "organizations"},
							{"localField", "location.id"},
							{"foreignField", "_id"},
							{"as", "organization"},
						}},
					},
					bson.D{
						{"$unwind", "$organization"},
					},
					bson.D{
						{"$addFields", bson.D{
							{"location.name", "$organization.name"},
						}},
					},
					bson.D{
						{"$project", bson.M{
							"organization": 0,
						}},
					},
				}},
				{"housingDonations", bson.A{
					bson.D{
						{"$match", bson.M{"location.type": utils.HousingLocationType}},
					},
					bson.D{
						{"$lookup", bson.D{
							{"from", "housings"},
							{"localField", "location.id"},
							{"foreignField", "_id"},
							{"as", "housing"},
						}},
					},
					bson.D{
						{"$unwind", "$housing"},
					},
					bson.D{
						{"$addFields", bson.D{
							{"location.name", "$housing.name"},
						}},
					},
					bson.D{
						{"$project", bson.M{
							"housing": 0,
						}},
					},
				}},
			}},
		},
		bson.D{
			{"$project", bson.D{
				{"donations", bson.D{
					{"$concatArrays", bson.A{"$organizationDonations", "$housingDonations"}},
				}},
			}},
		},
		bson.D{
			{"$unwind", "$donations"},
		},
		bson.D{
			{"$replaceRoot", bson.D{
				{"newRoot", "$donations"}},
			},
		},
	}

	cursor, err := repository.collection.Aggregate(context.TODO(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.TODO())

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}
