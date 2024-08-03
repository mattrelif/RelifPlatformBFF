package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/bff/entities"
	"relif/bff/utils"
	"time"
)

type UserPreferences struct {
	Language string `bson:"language,omitempty"`
	Timezone string `bson:"timezone,omitempty"`
}

type User struct {
	ID             string          `bson:"_id,omitempty"`
	FirstName      string          `bson:"first_name,omitempty"`
	LastName       string          `bson:"last_name,omitempty"`
	Email          string          `bson:"email,omitempty"`
	Password       string          `bson:"password,omitempty"`
	Phones         []string        `bson:"phones,omitempty"`
	OrganizationID string          `bson:"organization_id,omitempty"`
	Role           string          `bson:"role,omitempty"`
	PlatformRole   string          `bson:"platform_role,omitempty"`
	Status         string          `bson:"status,omitempty"`
	Preferences    UserPreferences `bson:"preferences,omitempty"`
	CreatedAt      time.Time       `bson:"created_at,omitempty"`
	UpdatedAt      time.Time       `bson:"updated_at,omitempty"`
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
		Preferences: entities.UserPreferences{
			Language: u.Preferences.Language,
			Timezone: u.Preferences.Timezone,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func NewUser(entity entities.User) User {
	return User{
		ID:           primitive.NewObjectID().Hex(),
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		Email:        entity.Email,
		Password:     entity.Password,
		Phones:       entity.Phones,
		Role:         entity.Role,
		PlatformRole: entity.PlatformRole,
		Status:       utils.ActiveStatus,
		Preferences: UserPreferences{
			Language: entity.Preferences.Language,
			Timezone: entity.Preferences.Timezone,
		},
		CreatedAt: time.Now(),
	}
}

func NewUpdatedUser(entity entities.User) User {
	return User{
		FirstName:    entity.FirstName,
		LastName:     entity.LastName,
		Email:        entity.Email,
		Password:     entity.Password,
		Phones:       entity.Phones,
		Role:         entity.Role,
		PlatformRole: entity.PlatformRole,
		Status:       entity.Status,
		Preferences: UserPreferences{
			Language: entity.Preferences.Language,
			Timezone: entity.Preferences.Timezone,
		},
		UpdatedAt: time.Now(),
	}
}
