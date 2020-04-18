package main

import (
	"context"
	"counter-app/pkg/app"
	"counter-app/pkg/webserver"
	"log"
	"os"
	"os/signal"
)

// TODO add populate access log
// TODO add rate limiter,
// TODO make it persistent
// TODO scale it with multiple instances
// TODO use flags to add params
// TODO performance tests
// TODO fix swagger, maybe CORS problem

func main() {
	log.Println("--- start ---")
	ctx, cancel := context.WithCancel(context.Background())

	application := app.New(ctx)
	server := webserver.New(application, "3000")

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	log.Println("Press CTRL + c to graceful shutdown ...")
	<-sigint

	log.Println("Shutting down ...")
	cancel()
	server.Stop()
	application.GracefulShutDown()

	log.Println(" --- end ---")
}
