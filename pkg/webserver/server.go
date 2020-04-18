package webserver

import (
	"context"
	"counter-app/pkg/app"
	"log"
	"net/http"
)

type Server struct {
	app        *app.App
	port       string
	httpServer http.Server
}

func New(app *app.App, port string) Server {
	router := http.NewServeMux()
	router.Handle("/keywords", accessLog(auth(final(app))))
	router.Handle("/swagger", accessLog(auth(swaggerInfo())))

	s := Server{
		app:  app,
		port: port,
		httpServer: http.Server{
			Addr:           ":" + port,
			Handler:        router,
			MaxHeaderBytes: 10000,
		},
	}
	go s.start()
	return s
}

func (s *Server) start() {
	log.Printf("WebServer starts and listen on port: %s", s.port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

func (s *Server) Stop() {
	log.Println("WebServer stopping ...")
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
		return
	}
	log.Println("WebServer stops with success.")
}
