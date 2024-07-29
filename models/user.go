package models

import (
	"fmt"
	"relif/bff/entities"
	"time"
)

type UserPreferences struct {
	Language string `bson:"language"`
	Timezone string `bson:"timezone"`
}

type User struct {
	ID             string          `bson:"_id"`
	FirstName      string          `bson:"first_name"`
	LastName       string          `bson:"last_name"`
	Email          string          `bson:"email"`
	Password       string          `bson:"password"`
	Phones         []string        `bson:"phones"`
	OrganizationID string          `bson:"organization_id"`
	Role           string          `bson:"role"`
	PlatformRole   string          `bson:"platform_role"`
	Status         string          `bson:"status"`
	Country        string          `bson:"country"`
	Preferences    UserPreferences `bson:"preferences"`
	CreatedAt      time.Time       `bson:"created_at"`
	UpdatedAt      time.Time       `bson:"updated_at"`
}

func (u *User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func (u *User) ToEntity() entities.User {
	return entities.User{
		ID:             u.ID,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		FullName:       u.FullName(),
		Email:          u.Email,
		Password:       u.Password,
		Phones:         u.Phones,
		OrganizationID: u.OrganizationID,
		Role:           u.Role,
		PlatformRole:   u.PlatformRole,
		Status:         u.Status,
		Country:        u.Country,
		Preferences: entities.UserPreferences{
			Language: u.Preferences.Language,
			Timezone: u.Preferences.Timezone,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
