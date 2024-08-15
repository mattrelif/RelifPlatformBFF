package entities

type StorageRecord struct {
	ID            string
	Location      Location
	ProductTypeID string
	ProductType   ProductType
	Quantity      int
}
