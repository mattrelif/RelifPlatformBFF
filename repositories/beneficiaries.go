package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type Beneficiaries interface {
	Create(data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingId(housingId string) ([]entities.Beneficiary, error)
	FindManyByRoomId(roomId string) ([]entities.Beneficiary, error)
	FindManyByOrganizationId(organizationId string) ([]entities.Beneficiary, error)
	FindOneById(id string) (entities.Beneficiary, error)
	CountByEmail(email string) (int64, error)
	UpdateOneById(id string, data entities.Beneficiary) error
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

func (repository *mongoBeneficiaries) FindManyByHousingId(housingId string) ([]entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{"current_housing_id": housingId}

	cursor, err := repository.collection.Find(context.Background(), filter)
	defer cursor.Close(context.Background())

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return entityList, nil
}

func (repository *mongoBeneficiaries) FindManyByRoomId(roomId string) ([]entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{"current_room_id": roomId}

	cursor, err := repository.collection.Find(context.Background(), filter)
	defer cursor.Close(context.Background())

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return entityList, nil
}

func (repository *mongoBeneficiaries) FindManyByOrganizationId(organizationId string) ([]entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{"current_organization_id": organizationId}

	cursor, err := repository.collection.Find(context.Background(), filter)
	defer cursor.Close(context.Background())

	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.Background(), &modelList); err != nil {
		return nil, err
	}

	for _, model := range modelList {
		entityList = append(entityList, model.ToEntity())
	}

	return entityList, nil
}

func (repository *mongoBeneficiaries) FindOneById(id string) (entities.Beneficiary, error) {
	var model models.Beneficiary

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Beneficiary{}, utils.ErrBeneficiaryNotFound
		}
		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) CountByEmail(email string) (int64, error) {
	filter := bson.M{"email": email}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *mongoBeneficiaries) UpdateOneById(id string, data entities.Beneficiary) error {
	model := models.NewUpdatedBeneficiary(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
