package main

import (
	"context"
	"counter-app/pkg/app"
	"counter-app/pkg/webserver"
	"log"
	"os"
	"os/signal"
)

// TODO make it thread safe
// TODO scale it with multiple instances
// TODO use flags to add params
// TODO improve disk tests
// TODO performance tests
// TODO add rate limiter, limit number of concurrent requests
// TODO fix swagger, maybe CORS problem
// TODO add details on access log

func main() {
	log.Println("--- start ---")
	ctx, cancel := context.WithCancel(context.Background())
	application, err := app.New(ctx, app.AppConf{})
	if err != nil {
		log.Fatal(err)
	}

	server := webserver.New(application, "3000")

	// TODO fi here, webserver starts after this, see console messages
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	log.Println("[Info] Press CTRL + c to graceful shutdown ...")
	<-sigint

	log.Println("[Info] Shutting down ...")
	cancel()
	server.Stop()
	application.GracefulShutDown()

	log.Println(" --- end ---")
}
