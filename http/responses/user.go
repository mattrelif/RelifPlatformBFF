package responses

import (
	"relif/bff/entities"
	"time"
)

type Users []User

type UserPreferences struct {
	Language string `json:"language"`
	Timezone string `json:"timezone"`
}

type User struct {
	ID             string          `json:"id"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	FullName       string          `json:"full_name"`
	Email          string          `json:"email"`
	Phones         []string        `json:"phones"`
	OrganizationID string          `json:"organization_id"`
	Role           string          `json:"role"`
	PlatformRole   string          `json:"platform_role"`
	Status         string          `json:"status"`
	Preferences    UserPreferences `json:"preferences"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

func NewUser(entity entities.User) User {
	return User{
		ID:             entity.ID,
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		FullName:       entity.FullName,
		Email:          entity.Email,
		Phones:         entity.Phones,
		OrganizationID: entity.OrganizationID,
		Role:           entity.Role,
		PlatformRole:   entity.PlatformRole,
		Status:         entity.Status,
		Preferences: UserPreferences{
			Language: entity.Preferences.Language,
			Timezone: entity.Preferences.Timezone,
		},
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func NewUsers(entityList []entities.User) Users {
	var response Users

	for _, entity := range entityList {
		response = append(response, NewUser(entity))
	}

	return response
}
