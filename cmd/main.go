package main

import (
	"context"
	"counter-app/pkg/memory"
	"counter-app/pkg/persistence/localdisk"
	"counter-app/pkg/publicapi"
	"counter-app/pkg/syncapi"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

// TODO add sync service
// TODO scale it with multiple instances
// TODO improve disk tests
// TODO performance tests
// TODO add rate limiter, limit number of concurrent requests
// TODO fix swagger, maybe CORS problem
// TODO add details on access log

// Default values used
type config struct {
	// have to be unique per node
	// env: COUNTER_NODE_ID
	nodeId int
	// represents endpoint where other nodes can be founded
	// env COUNTER_SEEDS (domain:syncport,domain:syncport...)
	seeds []string
	// where data will be stored on disk, can be absolute or relative
	// if relative, will create dir on execution place
	localStoragePath string
	localStorageFile string
	// listen and serve GET and POST on /keywords
	// env: COUNTER_PUBLIC_PORT
	publicPort string
	// used internal for syncapi with other nodes
	// env: COUNTER_SYNC_PORT
	syncPort                   string
	persistenceFlushIntervalMS int
}

func main() {
	logger := log.New(os.Stdout, "", 0)
	ctx, cancel := context.WithCancel(context.Background())

	conf := &config{}
	conf.localStoragePath = "data"
	conf.localStorageFile = "base"
	conf.persistenceFlushIntervalMS = 1000

	// parse env vars
	err := parseEnvVars(conf)
	if err != nil {
		logger.Fatal(err)
	}

	// Start memory
	memConf := &memory.Config{ID: conf.nodeId}
	mem := memory.NewMemory(ctx, logger, memConf)

	// Start persistence
	persistenceConf := &localdisk.Conf{
		Path:            conf.localStoragePath + strconv.Itoa(conf.nodeId),
		File:            conf.localStorageFile,
		FlushIntervalMS: conf.persistenceFlushIntervalMS,
	}
	persistence, err := localdisk.NewLocalDisk(ctx, logger, persistenceConf, mem)
	if err != nil {
		log.Fatal(err)
	}
	err = persistence.LoadVolatileStorage()
	if err != nil {
		log.Fatal(err)
	}

	logger.Println("[Info] Data loaded into volatile storage with success.")
	persistence.StartFlashing()
	logger.Printf("[Info] Start flashing into persistent storage periodically at: %d ms\n",
		conf.persistenceFlushIntervalMS,
	)

	// start public endpoint
	publicConf := &publicapi.Config{
		Port: conf.publicPort,
	}
	publicApi := publicapi.New(ctx, logger, publicConf, mem)
	publicApi.Start()

	// start sync endpoint
	syncConf := &syncapi.Config{
		Port: conf.syncPort,
	}
	syncApi := syncapi.New(ctx, logger, syncConf, mem)
	syncApi.Start()

	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	log.Println("[Info] Press CTRL + c to graceful shutdown ...")
	<-sigint

	log.Println("[Info] Shutting down ...")
	cancel()

	var message string
	// Stop public server
	message = fmt.Sprint("[Success] Stop HTTP server public api.")
	err = publicApi.Stop()
	if err != nil {
		message = fmt.Sprintf("[Error] HTTP server public api: %v", err)
	}
	logger.Println(message)

	message = fmt.Sprint("[Success] Stop flushing data to persistent storage.")
	err = persistence.StopFlashing()
	if err != nil {
		message = fmt.Sprintf("[Error] Flush to disk before close with error: %s", err)
		// TODO implement retry here, it is the last flush and it is important
	}
	logger.Println(message)

	// stop sync service
	// Stop public server
	message = fmt.Sprint("[Success] Stop HTTP sync public api.")
	err = syncApi.Stop()
	if err != nil {
		message = fmt.Sprintf("[Error] HTTP server sync api: %v", err)
	}
	logger.Println(message)

	log.Println(" --- End ---")
}

func parseEnvVars(conf *config) error {
	nodeId, err := strconv.Atoi(os.Getenv("COUNTER_NODE_ID"))
	if err != nil {
		return errors.New(
			fmt.Sprintf("[Error] %s Set on your environment an unique ID"+
				" integer for this node on \"COUNTER_NODE_ID\" var"+
				"nodes Ids should be incremental and starts with 0", err,
			),
		)
	}
	conf.nodeId = nodeId

	publicPort := os.Getenv("COUNTER_PUBLIC_PORT")
	if publicPort == "" {
		return errors.New(
			fmt.Sprintf("[Error] %s Set on your environment public fort"+
				" for this node on \"COUNTER_PUBLIC_PORT\" var", err,
			),
		)
	}
	conf.publicPort = publicPort

	syncPort := os.Getenv("COUNTER_SYNC_PORT")
	if syncPort == "" {
		return errors.New(
			fmt.Sprintf("[Error] %s Set on your environment public fort"+
				" for this node on \"COUNTER_SYNC_PORT\" var", err,
			),
		)
	}
	conf.syncPort = syncPort
	conf.seeds = strings.Fields(os.Getenv("COUNTER_SEEDS"))

	return nil
}
