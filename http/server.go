package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func NewServer(router http.Handler, port string) *Server {
	return &Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", port),
			Handler:      router,
			ReadTimeout:  time.Second * 10,
			WriteTimeout: time.Second * 10,
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
