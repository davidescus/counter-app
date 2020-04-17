package app

import (
	"counter-app/pkg/storage"
	"strings"
)

// Occurrence represents the number of occurrences
// for a specific string keyword
type Occurrence struct {
	Key   string
	Count uint64
}

type App struct {
	m *storage.Memory
}

func New() App {
	return App{m: storage.NewMemory()}
}

func (a *App) Start() {
	// TODO trigger command to load from disk
	// TODO start goroutine to sync to disc
	// TODO start goroutine to sync with another apps
}

// TODO add tests
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

// TODO add tests
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
