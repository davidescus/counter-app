package memory

import (
	"context"
	"hash/crc64"
	"log"
	"sync"
)

// TODO make it thread safe (mutex)
// TODO implement buckets to increase performance (basic with modulo)

type Config struct {
	// each node should have his specific id
	ID int
}

// Memory is the basic data structure that holds
// occurrences number for each keyword
type Memory struct {
	ctx         context.Context
	logger      *log.Logger
	id          int
	mu          *sync.RWMutex
	occurrences map[uint64]counts
	crcTable    *crc64.Table
}

type counts struct {
	// totals represents a slice of nodes, each node identifier
	// (starts with 0) will represent index in slice
	totals []uint64
}

func NewMemory(ctx context.Context, logger *log.Logger, conf *Config) *Memory {
	return &Memory{
		ctx:         ctx,
		logger:      logger,
		id:          conf.ID,
		mu:          &sync.RWMutex{},
		occurrences: make(map[uint64]counts),
		crcTable:    crc64.MakeTable(crc64.ISO),
	}
}

func (m *Memory) Increment(keyword []byte) {
	hash := m.generateHash(keyword)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.occurrences[hash]; !ok {
		m.occurrences[hash] = counts{
			totals: make([]uint64, m.id+1),
		}
	}

	if m.id+1 > len(m.occurrences[hash].totals) {
		counts := m.occurrences[hash]
		counts.totals = increaseSliceWithDefaultValues(counts.totals, m.id)
		m.occurrences[hash] = counts
	}

	m.occurrences[hash].totals[m.id]++
}

func (m *Memory) Get(keyword []byte) uint64 {
	hash := m.generateHash(keyword)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.occurrences[hash]; !ok {
		return uint64(0)
	}

	var total uint64
	for _, v := range m.occurrences[hash].totals {
		total += v
	}

	return total
}

// Import will load data into memory
func (m *Memory) Import(data map[uint64][]uint64) {
	for hash, totals := range data {
		m.occurrences[hash] = counts{
			totals: totals,
		}
	}
}

// Export will export all memory data
func (m *Memory) Export() map[uint64][]uint64 {
	export := make(map[uint64][]uint64)
	for hash, counts := range m.occurrences {
		export[hash] = counts.totals
	}
	return export
}

// by transforming []bytes into uint64 wil help us
// to create multiple buckets using modulo
func (m *Memory) generateHash(keyword []byte) uint64 {
	return crc64.Checksum(keyword, m.crcTable)
}

func increaseSliceWithDefaultValues(s []uint64, n int) []uint64 {
	if n-len(s) > 0 {
		for i := len(s); i < n; i++ {
			s = append(s, uint64(0))
		}
	}
	return s
}
