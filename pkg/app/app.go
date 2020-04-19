package app

import (
	"context"
	"counter-app/pkg/storage"
	"log"
	"strings"
	"time"
)

// Those are the default values, that will be used
// you can use your custom values when start app
// TODO create
const (
	storagePath         = "data"
	diskFlushIntervalMS = 1000
	syncIntervalMS      = 1000
)

// Occurrence represents the number of occurrences
// for a specific string keyword
type Occurrence struct {
	Key   string `json:"key"`
	Count uint64 `json:"count"`
}

// AppConf it is used for custom app configuration
type AppConf struct {
	StoragePath         string
	DiskFlushIntervalMS uint
	SyncIntervalMS      uint
}

type App struct {
	m                   *storage.Memory
	ctx                 context.Context
	diskFlushIntervalMS uint
	syncIntervalMS      uint
	stopSig             chan struct{}
	storagePath         string
	persistentStorage   *storage.Disk
}

// New will return an pointer to app ready to use
func New(ctx context.Context, conf AppConf) (*App, error) {
	app := App{
		m:       storage.NewMemory(),
		ctx:     ctx,
		stopSig: make(chan struct{}, 1),
	}
	app.applyDefaultsIfNeed(conf)

	return &app, app.load()
}

func (a *App) load() error {
	disk, err := storage.NewDisk(a.storagePath)
	if err != nil {
		return err
	}
	log.Println("[Info] memory warm up start")
	if err := disk.Load(a.m); err != nil {
		return err
	}
	log.Println("[Info] memory warm up finished")
	a.persistentStorage = disk
	go a.SyncAndFlush()

	return nil
}

func (a *App) SyncAndFlush() {
	log.Printf("[Success] starts flush and sync at each: %d ms", a.syncIntervalMS)
	for {
		select {
		case <-a.ctx.Done():
			// ugly, but needed, we do not want to lose data
			// during flush or sync ctx could be closed,
			// if we get a snapshot of data before flush or sync
			// will lose new entries
			a.Flush()
			a.Sync()
			a.stopSig <- struct{}{}
			return
		default:
			a.Flush()
			a.Sync()
			time.Sleep(time.Duration(a.syncIntervalMS) * time.Millisecond)
		}
	}
}

func (a *App) GracefulShutDown() {
	log.Println("Stopping app ...")
	<-a.stopSig
	log.Println("App stopped with success ...")
}

func (a *App) Flush() {
	err := a.persistentStorage.Flush(a.m)
	// TODO implement retry IMPORTANT
	if err != nil {
		log.Println("[Error] ", err)
	}

}

func (a *App) Sync() {
	//log.Println("Syncing with others app")
	// TODO implement this
	time.Sleep(300 * time.Millisecond)
	//log.Println("Finish to sync with others app")
}

// Store will transform string to lowercase,
// split it in keywords and store it in memory
func (a *App) Store(text string) error {
	// TODO we write into memory, so errors could raise when not enough space
	// required to use lowercase
	keywords := splitInKeywords(strings.ToLower(text))
	for _, keyword := range keywords {
		a.m.Increment(keyword)
	}
	return nil
}

// Get will return a collection of keys
// with their associated count occurrence number
func (a *App) Get(keys []string) []Occurrence {
	var occurrences []Occurrence
	for _, key := range keys {
		// required to use lowercase
		count := a.m.Get([]byte(strings.ToLower(key)))
		occurrences = append(occurrences, Occurrence{
			Key:   key,
			Count: count,
		})
	}
	return occurrences
}

func splitInKeywords(s string) [][]byte {
	var keywords [][]byte
	for _, key := range strings.Fields(s) {
		keywords = append(keywords, []byte(key))
	}
	return keywords
}

// if values are not sets on AppConf will apply default values
func (a *App) applyDefaultsIfNeed(conf AppConf) {
	if a.storagePath = conf.StoragePath; a.storagePath == "" {
		a.storagePath = storagePath
	}
	if a.diskFlushIntervalMS = conf.DiskFlushIntervalMS; a.diskFlushIntervalMS == 0 {
		a.diskFlushIntervalMS = diskFlushIntervalMS
	}
	if a.syncIntervalMS = conf.SyncIntervalMS; a.syncIntervalMS == 0 {
		a.syncIntervalMS = syncIntervalMS
	}
}
