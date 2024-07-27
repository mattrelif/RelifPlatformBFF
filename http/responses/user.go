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
	Email          string          `json:"email"`
	Phones         []string        `json:"phones"`
	OrganizationID string          `json:"organization_id"`
	Role           string          `json:"role"`
	PlatformRoleID string          `json:"platform_role_id"`
	Status         string          `json:"status"`
	Country        string          `json:"country"`
	Preferences    UserPreferences `json:"preferences"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	LastActivityAt time.Time       `json:"last_activity_at"`
}

func NewUser(entity entities.User) User {
	return User{
		ID:             entity.ID,
		FirstName:      entity.FirstName,
		LastName:       entity.LastName,
		Email:          entity.Email,
		Phones:         entity.Phones,
		OrganizationID: entity.OrganizationID,
		Role:           entity.Role,
		PlatformRoleID: entity.PlatformRoleID,
		Status:         entity.Status,
		Country:        entity.Country,
		Preferences: UserPreferences{
			Language: entity.Preferences.Language,
			Timezone: entity.Preferences.Timezone,
		},
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		LastActivityAt: entity.LastActivityAt,
	}
}

func NewUsers(entityList []entities.User) Users {
	var response Users

	for _, entity := range entityList {
		response = append(response, NewUser(entity))
	}

	return response
}
