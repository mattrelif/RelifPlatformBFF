package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"relif/platform-bff/utils"
)

type Server struct {
	server *http.Server
}

func NewServer(router http.Handler, port string, readTimeout, writeTimeout utils.Duration) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      router,
			ReadTimeout:  readTimeout.Duration,
			WriteTimeout: writeTimeout.Duration,
		},
	}
}

func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}
