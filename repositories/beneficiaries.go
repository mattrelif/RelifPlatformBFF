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

type Beneficiaries interface {
	Create(data entities.Beneficiary) (entities.Beneficiary, error)
	FindManyByHousingId(housingId string) ([]entities.Beneficiary, error)
	FindManyByRoomId(roomId string) ([]entities.Beneficiary, error)
	FindOneById(id string) (entities.Beneficiary, error)
	FindOneAndUpdateById(id string, data entities.Beneficiary) (entities.Beneficiary, error)
	UpdateOneById(id string, data entities.Beneficiary) error
	DeleteOneById(id string) error
}

type mongoBeneficiaries struct {
	collection *mongo.Collection
}

func NewMongoBeneficiares(database *mongo.Database) Beneficiaries {
	return &mongoBeneficiaries{
		collection: database.Collection("beneficiaries"),
	}
}

func (repository *mongoBeneficiaries) Create(data entities.Beneficiary) (entities.Beneficiary, error) {
	emergencyContacts := make([]models.EmergencyContact, 0)

	for _, contact := range data.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, models.EmergencyContact{Relationship: contact.Relationship, FullName: contact.FullName, Emails: contact.Emails, Phones: contact.Phones})
	}

	model := models.Beneficiary{
		ID:       primitive.NewObjectID().Hex(),
		FullName: data.FullName,
		Email:    data.Email,
		Document: models.Document{
			Type:  data.Document.Type,
			Value: data.Document.Value,
		},
		Birthdate:       data.Birthdate,
		Phones:          data.Phones,
		CivilStatus:     data.CivilStatus,
		SpokenLanguages: data.SpokenLanguages,
		Education:       data.Education,
		Address: models.Address{
			StreetNumber: data.Address.StreetNumber,
			StreetName:   data.Address.StreetName,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Status: data.Status,
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

func (repository *mongoBeneficiaries) FindOneById(id string) (entities.Beneficiary, error) {
	var model models.Beneficiary

	filter := bson.M{"_id": id}

	if err := repository.collection.FindOne(context.Background(), filter).Decode(&model); err != nil {
		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) FindOneAndUpdateById(id string, data entities.Beneficiary) (entities.Beneficiary, error) {
	emergencyContacts := make([]models.EmergencyContact, 0)

	for _, contact := range data.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, models.EmergencyContact{Relationship: contact.Relationship, FullName: contact.FullName, Emails: contact.Emails, Phones: contact.Phones})
	}

	model := models.Beneficiary{
		ID:       primitive.NewObjectID().Hex(),
		FullName: data.FullName,
		Email:    data.Email,
		Document: models.Document{
			Type:  data.Document.Type,
			Value: data.Document.Value,
		},
		Birthdate:       data.Birthdate,
		Phones:          data.Phones,
		CivilStatus:     data.CivilStatus,
		SpokenLanguages: data.SpokenLanguages,
		Education:       data.Education,
		Address: models.Address{
			StreetNumber: data.Address.StreetNumber,
			StreetName:   data.Address.StreetName,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Status: data.Status,
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

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&model); err != nil {
		return entities.Beneficiary{}, err
	}

	return model.ToEntity(), nil
}

func (repository *mongoBeneficiaries) UpdateOneById(id string, data entities.Beneficiary) error {
	emergencyContacts := make([]models.EmergencyContact, 0)

	for _, contact := range data.EmergencyContacts {
		emergencyContacts = append(emergencyContacts, models.EmergencyContact{Relationship: contact.Relationship, FullName: contact.FullName, Emails: contact.Emails, Phones: contact.Phones})
	}

	model := models.Beneficiary{
		ID:       primitive.NewObjectID().Hex(),
		FullName: data.FullName,
		Email:    data.Email,
		Document: models.Document{
			Type:  data.Document.Type,
			Value: data.Document.Value,
		},
		Birthdate:       data.Birthdate,
		Phones:          data.Phones,
		CivilStatus:     data.CivilStatus,
		SpokenLanguages: data.SpokenLanguages,
		Education:       data.Education,
		Address: models.Address{
			StreetNumber: data.Address.StreetNumber,
			StreetName:   data.Address.StreetName,
			ZipCode:      data.Address.ZipCode,
			District:     data.Address.District,
			City:         data.Address.City,
			Country:      data.Address.Country,
		},
		Status: data.Status,
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

	if err := repository.collection.FindOneAndUpdate(context.Background(), filter, update).Err(); err != nil {
		return err
	}

	return nil
}

func (repository *mongoBeneficiaries) DeleteOneById(id string) error {
	filter := bson.M{"_id": id}

	if err := repository.collection.FindOneAndDelete(context.Background(), filter).Err(); err != nil {
		return err
	}

	return nil
}
