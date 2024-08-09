package models

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"relif/platform-bff/entities"
	"relif/platform-bff/utils"
	"time"
)

type UserPreferences struct {
	Language string `bson:"language,omitempty"`
	Timezone string `bson:"timezone,omitempty"`
}

func (preferences *UserPreferences) ToEntity() entities.UserPreferences {
	return entities.UserPreferences{
		Language: preferences.Language,
		Timezone: preferences.Timezone,
	}
}

func NewUserPreferences(entity entities.UserPreferences) UserPreferences {
	return UserPreferences{
		Language: entity.Language,
		Timezone: entity.Timezone,
	}
}

type FindUser struct {
	ID             string          `bson:"_id,omitempty"`
	FirstName      string          `bson:"first_name,omitempty"`
	LastName       string          `bson:"last_name,omitempty"`
	Email          string          `bson:"email,omitempty"`
	Password       string          `bson:"password,omitempty"`
	Phones         []string        `bson:"phones,omitempty"`
	OrganizationID string          `bson:"organization_id,omitempty"`
	Organization   Organization    `bson:"organization,omitempty"`
	Role           string          `bson:"role,omitempty"`
	PlatformRole   string          `bson:"platform_role,omitempty"`
	Status         string          `bson:"status,omitempty"`
	Preferences    UserPreferences `bson:"preferences,omitempty"`
	CreatedAt      time.Time       `bson:"created_at,omitempty"`
	UpdatedAt      time.Time       `bson:"updated_at,omitempty"`
}

func (user *FindUser) ToEntity() entities.User {
	return entities.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		FullName:       fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:          user.Email,
		Password:       user.Password,
		Phones:         user.Phones,
		OrganizationID: user.OrganizationID,
		Organization:   user.Organization.ToEntity(),
		Role:           user.Role,
		PlatformRole:   user.PlatformRole,
		Status:         user.Status,
		Preferences:    user.Preferences.ToEntity(),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
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

func (user *User) ToEntity() entities.User {
	return entities.User{
		ID:             user.ID,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		FullName:       fmt.Sprintf("%s %s", user.FirstName, user.LastName),
		Email:          user.Email,
		Password:       user.Password,
		Phones:         user.Phones,
		OrganizationID: user.OrganizationID,
		Role:           user.Role,
		PlatformRole:   user.PlatformRole,
		Status:         user.Status,
		Preferences:    user.Preferences.ToEntity(),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func NewUser(entity entities.User) User {
	return User{
		ID:             primitive.NewObjectID().Hex(),
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		Email:          entity.Email,
		Password:       entity.Password,
		Phones:         entity.Phones,
		Role:           entity.Role,
		PlatformRole:   entity.PlatformRole,
		OrganizationID: entity.OrganizationID,
		Status:         utils.ActiveStatus,
		Preferences:    NewUserPreferences(entity.Preferences),
		CreatedAt:      time.Now(),
	}
}

func NewUpdatedUser(entity entities.User) User {
	return User{
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		Email:          entity.Email,
		Password:       entity.Password,
		Phones:         entity.Phones,
		Role:           entity.Role,
		PlatformRole:   entity.PlatformRole,
		OrganizationID: entity.OrganizationID,
		Status:         entity.Status,
		Preferences:    NewUserPreferences(entity.Preferences),
		UpdatedAt:      time.Now(),
	}
}
