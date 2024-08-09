package entities

type ProductCurrentStorage struct {
	ID             string
	HousingID      string
	Housing        Housing
	OrganizationID string
	Organization   Organization
	ProductTypeID  string
}
