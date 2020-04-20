package main

import (
	"context"
	"counter-app/pkg/memory"
	"counter-app/pkg/persistence/localdisk"
	"counter-app/pkg/publicapi"
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

// Default values used
const (
	localStoragePath           = "data"
	localStorageFile           = "base"
	publicAPIPort              = "3000"
	persistenceFlushIntervalMS = 1000
)

func main() {
	logger := log.New(os.Stdout, "", 0)
	ctx, cancel := context.WithCancel(context.Background())

	logger.Println("--- Start ---")

	// Start memory
	mem := memory.NewMemory()

	// Start persistence
	persistenceConf := &localdisk.Conf{
		Path:            localStoragePath,
		File:            localStorageFile,
		FlushIntervalMS: persistenceFlushIntervalMS,
	}
	persistence, err := localdisk.NewLocalDisk(ctx, logger, persistenceConf, mem)
	if err != nil {
		log.Fatal(err)
	}
	err = persistence.LoadVolatileStorage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[Info] Data loaded into volatile storage with success.")
	persistence.StartFlashing()
	log.Printf("[Info] Start flashing into persistent storage periodically at: %d ms\n",
		persistenceFlushIntervalMS,
	)

	// sync server

	// start Web client
	webConf := &publicapi.Config{
		Port: publicAPIPort,
	}
	publicApi := publicapi.New(ctx, logger, webConf, mem)
	publicApi.Start()

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	log.Println("[Info] Press CTRL + c to graceful shutdown ...")
	<-sigint

	log.Println("[Info] Shutting down ...")
	cancel()

	// Stop public API
	err = publicApi.Stop()
	if err != nil {
		log.Printf("[Error] HTTP server public api: %v", err)
	}
	log.Println("[Success] Stop HTTP server public api.")

	err = persistence.StopFlashing()
	// TODO implement retry here, it is the last flush and it is important
	logger.Println("[Error] Flush to disk before close with error: ", err)

	log.Println(" --- End ---")
}
