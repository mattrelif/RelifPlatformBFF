package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type Organizations []Organization

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Address     Address   `json:"address"`
	OwnerID     string    `json:"owner_id"`
	Status      string    `json:"status"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewOrganization(organization entities.Organization) Organization {
	return Organization{
		ID:          organization.ID,
		Name:        organization.Name,
		Description: organization.Description,
		Address:     NewAddress(organization.Address),
		OwnerID:     organization.OwnerID,
		Status:      organization.Status,
		Type:        organization.Type,
		CreatedAt:   organization.CreatedAt,
		UpdatedAt:   organization.UpdatedAt,
	}
}

func NewOrganizations(entityList []entities.Organization) Organizations {
	var response Organizations

	for _, entity := range entityList {
		response = append(response, NewOrganization(entity))
	}

	return response
}
