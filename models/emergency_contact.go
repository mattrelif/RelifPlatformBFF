package models

type EmergencyContact struct {
	Relationship string   `bson:"relationship"`
	FullName     string   `bson:"full_name"`
	Emails       []string `bson:"emails"`
	Phones       []string `bson:"phones"`
}
