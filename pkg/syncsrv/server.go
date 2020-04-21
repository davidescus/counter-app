package syncsrv

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server will serve sync api requests
type server struct {
	ctx     context.Context
	logger  *log.Logger
	storage Storage
	port    string
	http    http.Server
}

// will start sync endpoint
func (s *server) start() {
	router := http.NewServeMux()
	router.Handle("/sync", accessLog(auth(final(s.storage))))

	s.http = http.Server{
		Addr:           ":" + s.port,
		Handler:        router,
		MaxHeaderBytes: 10000,
	}

	go func() {
		s.logger.Printf("[Success] SyncService HTTP start on: %s", s.port)
		if err := s.http.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatalf("[Error] SyncService ListenAndServe: %v", err)
		}
	}()

	<-time.After(1 * time.Second)
}

func (s *server) stop() {
	if err := s.http.Shutdown(context.Background()); err != nil {
		s.logger.Printf("[Error] SyncService HTTP stop: %v", err)
		return
	}
	s.logger.Println("[Success] SyncService HTTP stop.")
}
