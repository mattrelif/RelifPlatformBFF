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

type VoluntaryPeople interface {
	Create(data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	FindManyByOrganizationId(organizationId, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error)
	FindOneById(id string) (entities.VoluntaryPerson, error)
	CountByEmail(email string) (int64, error)
	UpdateOneById(id string, data entities.VoluntaryPerson) error
}

type mongoVoluntaryPeople struct {
	collection *mongo.Collection
}

func NewMongoVoluntaryPeople(database *mongo.Database) VoluntaryPeople {
	return &mongoVoluntaryPeople{
		collection: database.Collection("voluntary_people"),
	}
}

func (repository *mongoVoluntaryPeople) Create(data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	model := models.NewVoluntaryPerson(data)

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoVoluntaryPeople) FindManyByOrganizationId(organizationId, search string, limit, offset int64) (int64, []entities.VoluntaryPerson, error) {
	var filter bson.M

	entityList := make([]entities.VoluntaryPerson, 0)
	modelsList := make([]models.VoluntaryPerson, 0)

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
					"full_name": bson.D{
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

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"full_name": 1})
	cursor, err := repository.collection.Find(context.Background(), filter, opts)

	if err != nil {
		return 0, nil, err
	}

	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &modelsList); err != nil {
		return 0, nil, err
	}

	for _, model := range modelsList {
		entityList = append(entityList, model.ToEntity())
	}

	return count, entityList, nil
}

func (repository *mongoVoluntaryPeople) FindOneById(id string) (entities.VoluntaryPerson, error) {
	var model models.VoluntaryPerson

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
			return entities.VoluntaryPerson{}, utils.ErrVoluntaryPersonNotFound
		}

		return entities.VoluntaryPerson{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoVoluntaryPeople) CountByEmail(email string) (int64, error) {
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

func (repository *mongoVoluntaryPeople) UpdateOneById(id string, data entities.VoluntaryPerson) error {
	model := models.NewUpdatedVoluntaryPerson(data)

	update := bson.M{"$set": &model}

	if _, err := repository.collection.UpdateByID(context.Background(), id, update); err != nil {
		return err
	}

	return nil
}
