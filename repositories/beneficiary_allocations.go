package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	modelList := make([]models.BeneficiaryAllocation, 0)
	entityList := make([]entities.BeneficiaryAllocation, 0)

	filter := bson.M{"beneficiary_id": beneficiaryID}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

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

func (repository *mongoBeneficiaryAllocations) FindManyByHousingIDPaginated(housingID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	modelList := make([]models.BeneficiaryAllocation, 0)
	entityList := make([]entities.BeneficiaryAllocation, 0)

	filter := bson.M{"$or": bson.A{bson.M{"old_housing_id": housingID}, bson.M{"housing_id": housingID}}}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

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

func (repository *mongoBeneficiaryAllocations) FindManyByRoomIDPaginated(roomID string, offset, limit int64) (int64, []entities.BeneficiaryAllocation, error) {
	modelList := make([]models.BeneficiaryAllocation, 0)
	entityList := make([]entities.BeneficiaryAllocation, 0)

	filter := bson.M{"$or": bson.A{bson.M{"old_room_id": roomID}, bson.M{"room_id": roomID}}}
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"created_at": -1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

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
