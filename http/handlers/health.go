package handlers

import "net/http"

type Health struct{}

func NewHealth() *Health {
	return &Health{}
}

func (handler *Health) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}
