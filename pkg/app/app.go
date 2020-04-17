package app

import (
	"counter-app/pkg/storage"
	"strings"
)

type Occurrence struct {
	Key   string
	Count int64
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
// TODO check if we need to return error if we write direct to ram memory
// Store will transform string to lowercase,
// split it in keywords and store it in memory
func (a *App) Store(input string) error {
	// required to use lowercase
	keywords := splitInKeywords(strings.ToLower(input))
	for _, keyword := range keywords {
		a.m.Set(keyword)
	}
	return nil
}

// TODO add tests
// Get will return a collection of keys
// with their associated count occurrence number
func (a *App) Get(input string) []Occurrence {
	var occurrences []Occurrence
	// required to use lowercase
	keywords := splitInKeywords(strings.ToLower(input))
	for _, keyword := range keywords {
		count := a.m.Get(keyword)
		occurrences = append(occurrences, Occurrence{
			Key:   string(keyword),
			Count: count,
		})
	}
	return occurrences
}

// TODO add tests
func splitInKeywords(s string) [][]byte {
	var keywords [][]byte
 	for _, key := range strings.Split(s, " ") {
		keywords = append(keywords, []byte(key))
	}
	return keywords
}
