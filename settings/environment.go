package settings

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Environment struct {
	AWS struct {
		Region string `default:"us-east-1"`
	}

	Router struct {
		Context string `default:"/api/v1"`
	}

	Server struct {
		Port         int           `default:"8080"`
		WriteTimeout time.Duration `default:"10s"`
		ReadTimeout  time.Duration `default:"10s"`
	}

	Mongo struct {
		URI               string        `default:"mongodb://127.0.0.1:27017"`
		Database          string        `default:"test"`
		ConnectionTimeout time.Duration `default:"10s"`
	}
}

func NewEnvironment() (*Environment, error) {
	var environment Environment

	if err := envconfig.Process("", &environment); err != nil {
		return nil, err
	}

	return &environment, nil
}
