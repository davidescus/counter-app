package syncapi

import (
	"context"
	"log"
	"net/http"
)

// Config for public api HTTP server
type Config struct {
	Port string
}

// Server will serve public api requests
type Server struct {
	ctx     context.Context
	logger  *log.Logger
	storage Storage
	port    string
	http    http.Server
}

// New instance of public api server
func New(ctx context.Context, logger *log.Logger, conf *Config, memory Storage) *Server {
	return &Server{
		ctx:     ctx,
		logger:  logger,
		storage: memory,
		port:    conf.Port,
	}
}

// Start will start public api HTTP server
func (s *Server) Start() {
	router := http.NewServeMux()
	router.Handle("/sync", accessLog(auth(final(s.storage))))

	s.http = http.Server{
		Addr:           ":" + s.port,
		Handler:        router,
		MaxHeaderBytes: 10000,
	}

	go func() {
		s.logger.Printf("[Success] Start HTTP sync api server on port: %s", s.port)
		if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatalf("[Error] HTTP sync server ListenAndServe: %v", err)
		}
	}()
}

// Stop will stop public api HTTP server
func (s *Server) Stop() error {
	return s.http.Shutdown(context.Background())
}