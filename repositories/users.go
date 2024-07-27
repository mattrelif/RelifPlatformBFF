package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"relif/bff/entities"
	"relif/bff/models"
)

type Users interface {
	CreateUser(user entities.User) (string, error)
	FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error)
	FindOneById(id string) (entities.User, error)
	FindOneByEmail(email string) (entities.User, error)
	FindOneAndUpdateById(id string, user entities.User) (entities.User, error)
	UpdateOneById(id string, user entities.User) error
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

func (repository *usersMongo) CreateUser(user entities.User) (string, error) {
	model := models.User{
		ID:             primitive.NewObjectID().Hex(),
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Password:       user.Password,
		Phones:         user.Phones,
		Role:           user.Role,
		PlatformRoleID: user.PlatformRoleID,
		Status:         user.Status,
		Country:        user.Country,
		Preferences: models.UserPreferences{
			Language: user.Preferences.Language,
			Timezone: user.Preferences.Timezone,
		},
		CreatedAt: user.CreatedAt,
	}

	result, err := repository.collection.InsertOne(context.Background(), &model)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (repository *usersMongo) FindManyByOrganizationId(organizationId string, offset, limit int64) (int64, []entities.User, error) {
	modelList := make([]models.User, 0)
	entityList := make([]entities.User, 0)

	oid, err := primitive.ObjectIDFromHex(organizationId)

	if err != nil {
		return 0, nil, err
	}

	filter := bson.M{"organization_id": oid}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	opts := options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"first_name": 1})
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

func (repository *usersMongo) FindOneById(id string) (entities.User, error) {
	var model models.User

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.User{}, err
	}

	filter := bson.M{"_id": oid}
	if err = repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
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

func (repository *usersMongo) FindOneAndUpdateById(id string, user entities.User) (entities.User, error) {
	model := models.User{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Password:       user.Password,
		Phones:         user.Phones,
		Role:           user.Role,
		PlatformRoleID: user.PlatformRoleID,
		Status:         user.Status,
		Country:        user.Country,
		Preferences: models.UserPreferences{
			Language: user.Preferences.Language,
			Timezone: user.Preferences.Timezone,
		},
		UpdatedAt:      user.UpdatedAt,
		LastActivityAt: user.LastActivityAt,
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return entities.User{}, err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.User{}, err
	}

	return model.ToEntity(), nil
}

func (repository *usersMongo) UpdateOneById(id string, user entities.User) error {
	model := models.User{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Password:       user.Password,
		Phones:         user.Phones,
		Role:           user.Role,
		PlatformRoleID: user.PlatformRoleID,
		Status:         user.Status,
		Country:        user.Country,
		Preferences: models.UserPreferences{
			Language: user.Preferences.Language,
			Timezone: user.Preferences.Timezone,
		},
		UpdatedAt:      user.UpdatedAt,
		LastActivityAt: user.LastActivityAt,
	}

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *usersMongo) DeleteOneById(id string) error {
	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err = repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
