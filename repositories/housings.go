package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type Housings interface {
	Create(data entities.Housing) (entities.Housing, error)
	FindManyByOrganizationID(organizationId, search string, limit, offset int64) (int64, []entities.Housing, error)
	FindOneByID(id string) (entities.Housing, error)
	UpdateOneById(id string, data entities.Housing) error
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

func (repository *mongoHousings) FindManyByOrganizationID(organizationId, search string, limit, offset int64) (int64, []entities.Housing, error) {
	var filter bson.M

	modelList := make([]models.FindHousing, 0)
	entityList := make([]entities.Housing, 0)

	if search != "" {
		filter = bson.M{
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
				{"let", bson.D{{"housingId", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$and", bson.A{
								bson.M{"$eq": bson.M{"$housing_id": "$$housingId"}},
								bson.M{"$ne": bson.M{"status": utils.InactiveStatus}},
							}},
						}},
					},
				}},
				{"as", "rooms"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "beneficiaries"},
				{"let", bson.D{{"housingId", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$and", bson.A{
								bson.M{"$eq": bson.M{"$current_housing_id": "$$housingId"}},
								bson.M{"$ne": bson.M{"status": utils.InactiveStatus}},
							}},
						}},
					},
				}},
				{"as", "beneficiaries"},
			}},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "$_id"},
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

func (repository *mongoHousings) FindOneByID(id string) (entities.Housing, error) {
	var model models.FindHousing

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

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"let", bson.D{{"housingId", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$and", bson.A{
								bson.M{"$eq": bson.M{"$housing_id": "$$housingId"}},
								bson.M{"$ne": bson.M{"status": utils.InactiveStatus}},
							}},
						}},
					},
				}},
				{"as", "rooms"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "beneficiaries"},
				{"let", bson.D{{"housingId", "$_id"}}},
				{"pipeline", bson.A{
					bson.D{
						{"$match", bson.D{
							{"$and", bson.A{
								bson.M{"$eq": bson.M{"$current_housing_id": "$$housingId"}},
								bson.M{"$ne": bson.M{"status": utils.InactiveStatus}},
							}},
						}},
					},
				}},
				{"as", "beneficiaries"},
			}},
		},
		bson.D{
			{"$group", bson.D{
				{"_id", "$_id"},
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
	defer cursor.Close(context.Background())

	if err != nil {
		return entities.Housing{}, err
	}

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.Housing{}, err
		}
	}

	return model.ToEntity(), nil
}

func (repository *mongoHousings) UpdateOneById(id string, data entities.Housing) error {
	model := models.NewUpdatedHousing(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
