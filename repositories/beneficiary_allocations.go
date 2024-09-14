package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/platform-bff/entities"
	"relif/platform-bff/models"
)

type BeneficiaryAllocations interface {
	Create(data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error)
	FindManyByBeneficiaryIDPaginated(beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByHousingIDPaginated(housingID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
	FindManyByRoomIDPaginated(roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error)
}

type mongoBeneficiaryAllocations struct {
	collection *mongo.Collection
}

func NewMongoBeneficiaryAllocations(database *mongo.Database) BeneficiaryAllocations {
	return &mongoBeneficiaryAllocations{
		collection: database.Collection("beneficiary_allocations"),
	}
}

func (repository *mongoBeneficiaryAllocations) Create(data entities.BeneficiaryAllocation) (entities.BeneficiaryAllocation, error) {
	model := models.NewBeneficiaryAllocation(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.BeneficiaryAllocation{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaryAllocations) FindManyByBeneficiaryIDPaginated(beneficiaryID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	entityList := make([]entities.BeneficiaryAllocation, 0)

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
				{"from", "housing_rooms"},
				{"localField", "old_room_id"},
				{"foreignField", "_id"},
				{"as", "old_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "old_housing_id"},
				{"foreignField", "_id"},
				{"as", "old_housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "room_id"},
				{"foreignField", "_id"},
				{"as", "room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "housing_id"},
				{"foreignField", "_id"},
				{"as", "housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.FindBeneficiaryAllocation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoBeneficiaryAllocations) FindManyByHousingIDPaginated(housingID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	entityList := make([]entities.BeneficiaryAllocation, 0)

	filter := bson.M{"$or": bson.A{bson.M{"old_housing_id": housingID}, bson.M{"housing_id": housingID}}}
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
				{"from", "housing_rooms"},
				{"localField", "old_room_id"},
				{"foreignField", "_id"},
				{"as", "old_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "old_housing_id"},
				{"foreignField", "_id"},
				{"as", "old_housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "room_id"},
				{"foreignField", "_id"},
				{"as", "room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "housing_id"},
				{"foreignField", "_id"},
				{"as", "housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.FindBeneficiaryAllocation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoBeneficiaryAllocations) FindManyByRoomIDPaginated(roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	entityList := make([]entities.BeneficiaryAllocation, 0)

	filter := bson.M{"$or": bson.A{bson.M{"old_room_id": roomID}, bson.M{"room_id": roomID}}}
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
				{"from", "housing_rooms"},
				{"localField", "old_room_id"},
				{"foreignField", "_id"},
				{"as", "old_room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "old_housing_id"},
				{"foreignField", "_id"},
				{"as", "old_housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$old_housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housing_rooms"},
				{"localField", "room_id"},
				{"foreignField", "_id"},
				{"as", "room"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$room"}, {"preserveNullAndEmptyArrays", true}}},
		},
		bson.D{
			{"$lookup", bson.D{
				{"from", "housings"},
				{"localField", "housing_id"},
				{"foreignField", "_id"},
				{"as", "housing"},
			}},
		},
		bson.D{
			{"$unwind", bson.D{{"path", "$housing"}, {"preserveNullAndEmptyArrays", true}}},
		},
	}

	cursor, err := repository.collection.Aggregate(context.Background(), pipeline)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var model models.FindBeneficiaryAllocation

		if err = cursor.Decode(&model); err != nil {
			return 0, nil, err
		}

		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}
