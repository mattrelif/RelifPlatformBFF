package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type ProductTypes []ProductType

type ProductType struct {
	ID             string       `json:"id"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	Brand          string       `json:"brand"`
	Category       string       `json:"category"`
	OrganizationID string       `json:"organization_id"`
	Organization   Organization `json:"organization"`
	UnitType       string       `json:"unit_type"`
	TotalInStorage int          `json:"total_in_storage"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func NewProductType(entity entities.ProductType) ProductType {
	return ProductType{
		ID:             entity.ID,
		Name:           entity.Name,
		Description:    entity.Description,
		Category:       entity.Category,
		OrganizationID: entity.OrganizationID,
		Brand:          entity.Brand,
		Organization:   NewOrganization(entity.Organization),
		UnitType:       entity.UnitType,
		TotalInStorage: entity.TotalInStorage,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

func NewProductTypes(entityList []entities.ProductType) ProductTypes {
	res := make([]ProductType, 0)

	for _, entity := range entityList {
		res = append(res, NewProductType(entity))
	}

	return res
}
