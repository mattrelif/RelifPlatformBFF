package repositories

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"relif/bff/utils"
)

type Beneficiaries interface {
	Create(data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.Beneficiary, error)
	FindManyByRoomId(roomId string, limit, offset int64) (int64, []entities.Beneficiary, error)
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.Beneficiary, error)
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

func (repository *mongoBeneficiaries) FindManyByHousingId(housingId string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"current_housing_id": housingId,
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

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"full_name": 1})

	cursor, err := repository.collection.Find(context.Background(), filter, opts)
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

func (repository *mongoBeneficiaries) FindManyByRoomId(roomId string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"current_room_id": roomId,
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
	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"full_name": 1})

	cursor, err := repository.collection.Find(context.Background(), filter, opts)
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

func (repository *mongoBeneficiaries) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.Beneficiary, error) {
	entityList := make([]entities.Beneficiary, 0)
	modelList := make([]models.Beneficiary, 0)

	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"current_organization_id": organizationId,
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

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"full_name": 1})

	cursor, err := repository.collection.Find(context.Background(), filter, opts)
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

func (repository *mongoBeneficiaries) FindOneById(id string) (entities.Beneficiary, error) {
	var model models.Beneficiary

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

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.Beneficiary{}, utils.ErrBeneficiaryNotFound
		}
		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) CountByEmail(email string) (int64, error) {
	filter := bson.M{
		"$and": bson.A{
			bson.M{
				"email": email,
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
