package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Health struct {
	mongoClient *mongo.Client
}

type HealthResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
	Time     string            `json:"timestamp"`
}

func NewHealth(mongoClient *mongo.Client) *Health {
	return &Health{
		mongoClient: mongoClient,
	}
}

func (handler *Health) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:   "healthy",
		Services: make(map[string]string),
		Time:     time.Now().UTC().Format(time.RFC3339),
	}

	// Check MongoDB connectivity
	if handler.mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := handler.mongoClient.Ping(ctx, readpref.Primary())
		if err != nil {
			response.Services["mongodb"] = "unhealthy: " + err.Error()
			response.Status = "degraded"
		} else {
			response.Services["mongodb"] = "healthy"
		}
	} else {
		response.Services["mongodb"] = "not configured"
		response.Status = "degraded"
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Set status code based on overall health
	if response.Status == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	// Return JSON response
	json.NewEncoder(w).Encode(response)
}

// Simple health check for load balancers (just returns "Healthy" text)
func (handler *Health) SimpleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}
