package app

import (
	"context"
	"counter-app/pkg/storage"
	"log"
	"strings"
	"time"
)

// Occurrence represents the number of occurrences
// for a specific string keyword
type Occurrence struct {
	Key   string `json:"key"`
	Count uint64 `json:"count"`
}

type App struct {
	m         *storage.Memory
	ctx       context.Context
	heartBeat time.Duration
	stopSig   chan struct{}
}

func New(ctx context.Context) *App {
	app := App{
		m:         storage.NewMemory(),
		ctx:       ctx,
		heartBeat: time.Duration(1000),
		stopSig:   make(chan struct{}, 1),
	}
	app.load()

	return &app
}

func (a *App) load() {
	// TODO trigger command to load from disk

	go a.SyncAndFlush()
}

func (a *App) SyncAndFlush() {
	log.Printf("Start flushing and synking data at each: %d milliseconds", a.heartBeat)
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
			time.Sleep(a.heartBeat * time.Millisecond)
		}
	}
}

func (a *App) GracefulShutDown() {
	log.Println("Stopping app ...")
	<-a.stopSig
	log.Println("App stopped with success ...")
}

func (a *App) Flush() {
	//log.Println("Flushing data to disk ...")
	// TODO implement this
	time.Sleep(200 * time.Millisecond)
	//log.Println("Finish flush to disk.")

}

func (a *App) Sync() {
	//log.Println("Syncing with others app")
	// TODO implement this
	time.Sleep(300 * time.Millisecond)
	//log.Println("Finish to sync with others app")
}

// TODO we write into memory, so errors could raise when not enough space
// Store will transform string to lowercase,
// split it in keywords and store it in memory
func (a *App) Store(text string) error {
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
