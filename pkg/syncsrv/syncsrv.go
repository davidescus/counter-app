package syncsrv

import (
	"context"
	"log"
	"time"
)

const syncURI = "/sync"

// Storage ...
type Storage interface {
	Merge(occurrences map[uint64][]uint64)
	Export() map[uint64][]uint64
}

type Occurrence struct {
	HashKeyword uint64   `json:"hash_keyword"`
	Totals      []uint64 `json:"totals"`
}

// Config for sync service
type Config struct {
	ListeningPort string
	Seeds         []string
	SyncInterval  time.Duration
}

// Sync will exchange data with other seedsj
type Sync struct {
	ctx                   context.Context
	logger                *log.Logger
	storage               Storage
	seeds                 []string
	syncInterval          time.Duration
	endSig, endSigSuccess chan struct{}
	server                server
	client                client
}

// New ...
func New(ctx context.Context, logger *log.Logger, conf *Config, memory Storage) *Sync {
	return &Sync{
		ctx:           ctx,
		logger:        logger,
		storage:       memory,
		seeds:         conf.Seeds,
		syncInterval:  conf.SyncInterval,
		endSig:        make(chan struct{}),
		endSigSuccess: make(chan struct{}),
		server: server{
			ctx:     ctx,
			logger:  logger,
			storage: memory,
			port:    conf.ListeningPort,
		},
		client: client{},
	}
}

func (s *Sync) Start() {
	s.server.start()
	s.startSync()
}

func (s *Sync) Stop() {
	s.endSig <- struct{}{}
	<-s.endSigSuccess
	s.server.stop()
	s.syncNow()
}

// will share data with others seeds periodically
func (s *Sync) startSync() {
	go func() {
		for {
			select {
			case <-s.ctx.Done():
			case <-s.endSig:
				s.endSigSuccess <- struct{}{}
				return
			default:
				s.syncNow()
				time.Sleep(s.syncInterval)
			}
		}
	}()
}

// on fist iteration will send all data to all seeds
// each time, will improve this by adding last updated
// on each key in storage
func (s *Sync) syncNow() {
	var occurrences []Occurrence
	for k, v := range s.storage.Export() {
		occurrences = append(occurrences, Occurrence{
			HashKeyword: k,
			Totals:      v,
		})
	}

	for _, seedAddr := range s.seeds {
		if err := s.client.send(occurrences, seedAddr+syncURI); err != nil {
			s.logger.Println("[Error} SyncNow: ", err)
		}
	}
}
