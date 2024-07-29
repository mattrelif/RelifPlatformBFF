package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
	"time"
)

type Users interface {
	CreateUser(data entities.User) (entities.User, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	FindOneAndUpdateById(id string, data entities.User) (entities.User, error)
	UpdateOneById(id string, data entities.User) error
	DeleteOneById(id string) error
}

type usersMongo struct {
	collection *mongo.Collection
}

func NewUsersMongo(database *mongo.Database) Users {
	return &usersMongo{
		collection: database.Collection("users"),
	}
}

func (repository *usersMongo) CreateUser(data entities.User) (entities.User, error) {
	model := models.User{
		ID:           primitive.NewObjectID().Hex(),
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		Password:     data.Password,
		Phones:       data.Phones,
		Role:         data.Role,
		PlatformRole: data.PlatformRole,
		Status:       data.Status,
		Country:      data.Country,
		Preferences: models.UserPreferences{
			Language: data.Preferences.Language,
			Timezone: data.Preferences.Timezone,
		},
		CreatedAt: time.Now(),
	}

	if _, err := repository.collection.InsertOne(context.Background(), &model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error) {
	modelList := make([]models.User, 0)
	entityList := make([]entities.User, 0)

	filter := bson.M{"organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"first_name": 1})
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

func (repository *usersMongo) FindOneById(id string) (entities.User, error) {
	var model models.User

	filter := bson.M{"_id": id}
	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) FindOneByEmail(email string) (entities.User, error) {
	var model models.User

	filter := bson.M{"email": email}
	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) FindOneAndUpdateById(id string, data entities.User) (entities.User, error) {
	model := models.User{
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		Password:     data.Password,
		Phones:       data.Phones,
		Role:         data.Role,
		PlatformRole: data.PlatformRole,
		Status:       data.Status,
		Country:      data.Country,
		Preferences: models.UserPreferences{
			Language: data.Preferences.Language,
			Timezone: data.Preferences.Timezone,
		},
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) UpdateOneById(id string, data entities.User) error {
	model := models.User{
		FirstName:    data.FirstName,
		LastName:     data.LastName,
		Email:        data.Email,
		Password:     data.Password,
		Phones:       data.Phones,
		Role:         data.Role,
		PlatformRole: data.PlatformRole,
		Status:       data.Status,
		Country:      data.Country,
		Preferences: models.UserPreferences{
			Language: data.Preferences.Language,
			Timezone: data.Preferences.Timezone,
		},
		UpdatedAt: time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *usersMongo) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}
	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
