package settings

import (
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"os"
	"strings"
	"time"
)

var ApplicationEnvironment = os.Getenv("APPLICATION_ENVIRONMENT")

type Environment struct {
	Email struct {
		Domain string `required:"true" json:"EMAIL_DOMAIN"`
	}

	AWS struct {
		Region string `default:"us-east-1" json:"AWS_REGION"`
	}

	Router struct {
		Context string `default:"/api/v1" json:"ROUTER_CONTEXT"`
	}

	Server struct {
		Port         int           `default:"8080" json:"SERVER_PORT"`
		WriteTimeout time.Duration `default:"10s" json:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `default:"10s" json:"SERVER_READ_TIMEOUT"`
	}

	Mongo struct {
		URI               string        `default:"mongodb://127.0.0.1:27017" json:"MONGO_URI"`
		Database          string        `default:"test" json:"MONGO_DATABASE"`
		ConnectionTimeout time.Duration `default:"10s" json:"MONGO_CONNECTION_TIMEOUT"`
	}
}

func NewEnvironment() (*Environment, error) {
	var environment Environment

	if strings.ToLower(ApplicationEnvironment) == "development" {
		if err := envconfig.Process("", &environment); err != nil {
			return nil, err
		}
	} else if strings.ToLower(ApplicationEnvironment) == "production" {
		variables := os.Getenv("ENVIRONMENT_VARIABLES")

		if err := json.Unmarshal([]byte(variables), &environment); err != nil {
			return nil, err
		}
	}

	return &environment, nil
}
