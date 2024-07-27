package repositories

import "go.mongodb.org/mongo-driver/mongo"

type Beneficiaries interface {
}

type mongoBeneficiaries struct {
	collection *mongo.Collection
}

func NewMongoBeneficiares(database *mongo.Database) Beneficiaries {
	return &mongoBeneficiaries{
		collection: database.Collection("beneficiaries"),
	}
}
