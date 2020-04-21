package main

import (
	"context"
	"counter-app/pkg/memory"
	"counter-app/pkg/persistence/localdisk"
	"counter-app/pkg/publicapi"
	"counter-app/pkg/syncsrv"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

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
	// used internal for syncsrv with other nodes
	// env: COUNTER_SYNC_PORT
	syncPort     string
	dumpInterval time.Duration
	syncInterval time.Duration
}

func main() {
	logger := log.New(os.Stdout, "", 0)
	ctx, cancel := context.WithCancel(context.Background())

	conf := &config{}
	dumpInt, _ := time.ParseDuration("1s")
	syncInt, _ := time.ParseDuration("2s")
	conf.localStoragePath = "data"
	conf.localStorageFile = "base"
	conf.dumpInterval = dumpInt
	conf.syncInterval = syncInt

	// parse env vars
	err := parseEnvVars(conf)
	if err != nil {
		logger.Fatal("[Fatal] ", err)
	}

	// StartPersistentDumps memory
	memConf := &memory.Config{ID: conf.nodeId}
	mem := memory.NewMemory(ctx, logger, memConf)

	// StartPersistentDumps persistence
	persistenceConf := &localdisk.Conf{
		Path:         conf.localStoragePath,
		File:         conf.localStorageFile + strconv.Itoa(conf.nodeId),
		DumpInterval: conf.dumpInterval,
	}
	persistence, err := localdisk.NewLocalDisk(ctx, logger, persistenceConf, mem)
	if err != nil {
		logger.Fatal("[Fatal] ", err)
	}
	err = persistence.LoadVolatileStorage()
	if err != nil {
		logger.Fatal("[Fatal] ", err)
	}
	logger.Println("[Info] Data loaded into volatile storage with success.")
	persistence.StartPersistentDumps()
	logger.Printf("[Info] StartPersistentDumps periodically at: %d ms\n",
		conf.dumpInterval,
	)

	// start public endpoint
	publicConf := &publicapi.Config{
		Port: conf.publicPort,
	}
	publicApi := publicapi.New(ctx, logger, publicConf, mem)
	publicApi.Start()

	// start sync service
	syncConf := &syncsrv.Config{
		Seeds:         conf.seeds,
		ListeningPort: conf.syncPort,
		SyncInterval:  conf.syncInterval,
	}
	sync := syncsrv.New(ctx, logger, syncConf, mem)
	sync.Start()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	log.Println("[Info] Press CTRL + c to graceful shutdown ...")

	// Wait, shutdown gracefully
	<-sigint
	logger.Println("[Info] Shutting down ...")
	cancel()

	var message string
	// Stop public server
	message = fmt.Sprint("[Success] PublicService HTTP stop.")
	err = publicApi.Stop()
	if err != nil {
		message = fmt.Sprintf("[Error] PublicService: %v", err)
	}
	logger.Println(message)

	// Stop persistence
	message = fmt.Sprint("[Success] PersistentDumps stop.")
	err = persistence.StopPersistenceDump()
	if err != nil {
		message = fmt.Sprintf("[Error] PersistentDump before close: %s", err)
	}
	logger.Println(message)

	// Stop sync service
	sync.Stop()
	logger.Println("[Success] SyncService stops.")

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
	conf.seeds = strings.Split(os.Getenv("COUNTER_SEEDS"), ",")

	return nil
}
