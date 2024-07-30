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

type VoluntaryPeople interface {
	Create(data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.VoluntaryPerson, error)
	FindOneById(id string) (entities.VoluntaryPerson, error)
	FindOneAndUpdateById(id string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error)
	DeleteOneById(id string) error
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
	emergencyContacts := make([]models.EmergencyContact, 0)

	for _, contact := range data.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, models.EmergencyContact{Relationship: contact.Relationship, FullName: contact.FullName, Emails: contact.Emails, Phones: contact.Phones})
	}

	model := models.VoluntaryPerson{
		ID:       primitive.NewObjectID().Hex(),
		FullName: data.FullName,
		Email:    data.Email,
		Document: models.Document{
			Type:  data.Document.Type,
			Value: data.Document.Value,
		},
		Birthdate: data.Birthdate,
		Phones:    data.Phones,
		Address: models.Address{
			StreetNumber: data.Address.StreetNumber,
			StreetName:   data.Address.StreetName,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Status:   data.Status,
		Segments: data.Segments,
		MedicalInformation: models.MedicalInformation{
			Allergies:                  data.MedicalInformation.Allergies,
			CurrentMedications:         data.MedicalInformation.CurrentMedications,
			RecurrentMedicalConditions: data.MedicalInformation.RecurrentMedicalConditions,
			HealthInsurancePlans:       data.MedicalInformation.HealthInsurancePlans,
			BloodType:                  data.MedicalInformation.BloodType,
			TakenVaccines:              data.MedicalInformation.TakenVaccines,
			MentalHealthHistory:        data.MedicalInformation.MentalHealthHistory,
			Height:                     data.MedicalInformation.Height,
			Weight:                     data.MedicalInformation.Weight,
			CigarettesUsage:            data.MedicalInformation.CigarettesUsage,
			AlcoholConsumption:         data.MedicalInformation.AlcoholConsumption,
			Disabilities:               data.MedicalInformation.Disabilities,
		},
		EmergencyContacts: emergencyContacts,
		CreatedAt:         time.Now(),
		Notes:             data.Notes,
	}

	if _, err := repository.collection.InsertOne(context.Background(), model); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoVoluntaryPeople) FindManyByOrganizationId(organizationId string, limit, offset int64) (int64, []entities.VoluntaryPerson, error) {
	entityList := make([]entities.VoluntaryPerson, 0)
	modelsList := make([]models.VoluntaryPerson, 0)

	filter := bson.M{"organization_id": organizationId}

	count, err := repository.collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, nil, err
	}

	opts := options.Find().SetLimit(limit).SetSkip(offset).SetSort(bson.M{"full_name": 1})

	cursor, err := repository.collection.Find(context.Background(), filter, opts)
	defer cursor.Close(context.Background())

	if err != nil {
		return 0, nil, err
	}

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

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoVoluntaryPeople) FindOneAndUpdateById(id string, data entities.VoluntaryPerson) (entities.VoluntaryPerson, error) {
	emergencyContacts := make([]models.EmergencyContact, 0)

	for _, contact := range data.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, models.EmergencyContact{Relationship: contact.Relationship, FullName: contact.FullName, Emails: contact.Emails, Phones: contact.Phones})
	}

	model := models.VoluntaryPerson{
		ID:       primitive.NewObjectID().Hex(),
		FullName: data.FullName,
		Email:    data.Email,
		Document: models.Document{
			Type:  data.Document.Type,
			Value: data.Document.Value,
		},
		Birthdate: data.Birthdate,
		Phones:    data.Phones,
		Address: models.Address{
			StreetNumber: data.Address.StreetNumber,
			StreetName:   data.Address.StreetName,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Status:   data.Status,
		Segments: data.Segments,
		MedicalInformation: models.MedicalInformation{
			Allergies:                  data.MedicalInformation.Allergies,
			CurrentMedications:         data.MedicalInformation.CurrentMedications,
			RecurrentMedicalConditions: data.MedicalInformation.RecurrentMedicalConditions,
			HealthInsurancePlans:       data.MedicalInformation.HealthInsurancePlans,
			BloodType:                  data.MedicalInformation.BloodType,
			TakenVaccines:              data.MedicalInformation.TakenVaccines,
			MentalHealthHistory:        data.MedicalInformation.MentalHealthHistory,
			Height:                     data.MedicalInformation.Height,
			Weight:                     data.MedicalInformation.Weight,
			CigarettesUsage:            data.MedicalInformation.CigarettesUsage,
			AlcoholConsumption:         data.MedicalInformation.AlcoholConsumption,
			Disabilities:               data.MedicalInformation.Disabilities,
		},
		EmergencyContacts: emergencyContacts,
		UpdatedAt:         time.Now(),
		Notes:             data.Notes,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": &model}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Err(); err != nil {
		return entities.VoluntaryPerson{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoVoluntaryPeople) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
