package responses

type FindMany[T any] struct {
	Count int64 `json:"count"`
	Data  T     `json:"data"`
}
