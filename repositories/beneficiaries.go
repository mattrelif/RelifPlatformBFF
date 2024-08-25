package repositories

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
	"relif/platform-bff/utils"
)

type Beneficiaries interface {
	Create(data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingIDPaginated(housingID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
	FindManyByRoomIDPaginated(roomID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
	FindManyByOrganizationIDPaginated(organizationID, search string, offset, limit int64) (int64, []entities.Beneficiary, error)
	FindOneByID(id string) (entities.Beneficiary, error)
	FindOneCompleteByID(id string) (entities.Beneficiary, error)
	CountByEmail(email string) (int64, error)
	UpdateOneByID(id string, data entities.Beneficiary) error
	DeleteOneByID(id string) error
}

type mongoBeneficiaries struct {
	collection *mongo.Collection
}

func NewMongoBeneficiaries(database *mongo.Database) Beneficiaries {
	return &mongoBeneficiaries{
		collection: database.Collection("beneficiaries"),
	}
}

func (repository *mongoBeneficiaries) Create(data entities.Beneficiary) (entities.Beneficiary, error) {
	model := models.NewBeneficiary(data)

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) FindManyByHousingIDPaginated(housingID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
	var filter bson.M

	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.FindBeneficiary, 0)

	if search != "" {
		filter = bson.M{
			"$and": bson.A{
				bson.M{
					"current_housing_id": housingID,
				},
				bson.M{
					"full_name": bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}
	} else {
		filter = bson.M{
			"current_housing_id": housingID,
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
			{"$sort", bson.M{"full_name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "current_organization_id"},
				{"foreignField", "_id"},
				{"as", "current_organization"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "current_housing_id"},
				{"foreignField", "_id"},
				{"as", "current_housing"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "current_room_id"},
				{"foreignField", "_id"},
				{"as", "current_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_room"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoBeneficiaries) FindManyByRoomIDPaginated(roomID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
	var filter bson.M

	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.FindBeneficiary, 0)

	if search != "" {
		filter = bson.M{
			"$and": bson.A{
				bson.M{
					"current_room_id": roomID,
				},
				bson.M{
					"full_name": bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}
	} else {
		filter = bson.M{
			"current_room_id": roomID,
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
			{"$sort", bson.M{"full_name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "current_organization_id"},
				{"foreignField", "_id"},
				{"as", "current_organization"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "current_housing_id"},
				{"foreignField", "_id"},
				{"as", "current_housing"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "current_room_id"},
				{"foreignField", "_id"},
				{"as", "current_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_room"}, {"preserveNullAndEmptyArrays", true}}},
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

func (repository *mongoBeneficiaries) FindManyByOrganizationIDPaginated(organizationID, search string, offset, limit int64) (int64, []entities.Beneficiary, error) {
	var filter bson.M

	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.FindBeneficiary, 0)

	if search != "" {
		filter = bson.M{
			"$and": bson.A{
				bson.M{
					"current_organization_id": organizationID,
				},
				bson.M{
					"full_name": bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}
	} else {
		filter = bson.M{
			"current_organization_id": organizationID,
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
			{"$sort", bson.M{"full_name": 1}},
		},
		bson.D{
			{"$skip", offset},
		},
		bson.D{
			{"$limit", limit},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "current_organization_id"},
				{"foreignField", "_id"},
				{"as", "current_organization"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "current_housing_id"},
				{"foreignField", "_id"},
				{"as", "current_housing"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "current_room_id"},
				{"foreignField", "_id"},
				{"as", "current_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_room"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var decoded map[string]interface{}

		if err = cursor.Decode(&decoded); err != nil {
			return 0, nil, err
		}

		fmt.Println(decoded)
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoBeneficiaries) FindOneByID(id string) (entities.Beneficiary, error) {
	var model models.Beneficiary

	filter := bson.M{
		"_id": id,
	}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Beneficiary{}, utils.ErrBeneficiaryNotFound
		}

		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) FindOneCompleteByID(id string) (entities.Beneficiary, error) {
	var model models.FindBeneficiary

	filter := bson.M{
		"_id": id,
	}

	pipeline := mongo.Pipeline{
		bson.D{
			{"$match", filter},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "organizations"},
				{"localField", "current_organization_id"},
				{"foreignField", "_id"},
				{"as", "current_organization"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "current_housing_id"},
				{"foreignField", "_id"},
				{"as", "current_housing"},
			}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "current_room_id"},
				{"foreignField", "_id"},
				{"as", "current_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_organization"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$current_room"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return entities.Beneficiary{}, err
	}

	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		if err = cursor.Decode(&model); err != nil {
			return entities.Beneficiary{}, err
		}
	} else {
		return entities.Beneficiary{}, utils.ErrBeneficiaryNotFound
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) CountByEmail(email string) (int64, error) {
	filter := bson.M{
		"email": email,
	}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *mongoBeneficiaries) UpdateOneByID(id string, data entities.Beneficiary) error {
	model := models.NewUpdatedBeneficiary(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}

func (repository *mongoBeneficiaries) DeleteOneByID(id string) error {
	filter := bson.M{"_id": id}

	if _, err := repository.collection.DeleteOne(context.Background(), filter); err != nil {
		return err
	}

	return nil
}
