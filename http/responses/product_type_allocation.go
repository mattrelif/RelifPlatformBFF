package responses

import (
	"relif/platform-bff/entities"
	"time"
)

type ProductTypeAllocation struct {
	ID             string    `json:"id,omitempty"`
	ProductTypeID  string    `json:"product_type_id,omitempty'"`
	Type           string    `json:"type,omitempty'"`
	From           Location  `json:"from,omitempty'"`
	To             Location  `json:"to,omitempty'"`
	OrganizationID string    `json:"organization_id,omitempty'"`
	Quantity       int       `json:"quantity,omitempty'"`
	CreatedAt      time.Time `json:"created_at,omitempty'"`
}

func NewProductTypeAllocation(entity entities.ProductTypeAllocation) ProductTypeAllocation {
	return ProductTypeAllocation{
		ID:             entity.ID,
		ProductTypeID:  entity.ProductTypeID,
		Type:           entity.Type,
		From:           NewLocation(entity.From),
		To:             NewLocation(entity.To),
		OrganizationID: entity.OrganizationID,
		Quantity:       entity.Quantity,
		CreatedAt:      entity.CreatedAt,
	}
}
