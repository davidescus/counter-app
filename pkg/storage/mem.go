package storage

import (
	"hash/crc64"
	"sync"
)

// TODO make it thread safe (mutex)
// TODO implement buckets to increase performance (basic with modulo)

// Memory is the basic data structure that holds
// number of occurrence for each keyword
type Memory struct {
	mu          *sync.RWMutex
	occurrences map[uint64]uint64
	crcTable    *crc64.Table
}

func NewMemory() *Memory {
	return &Memory{
		occurrences: make(map[uint64]uint64),
		crcTable:    crc64.MakeTable(crc64.ISO),
	}
}

func (m *Memory) Increment(keyword []byte) {
	m.occurrences[m.generateHash(keyword)]++
}

func (m *Memory) Get(keyword []byte) uint64 {
	return m.occurrences[m.generateHash(keyword)]
}

// by transforming []bytes into uint64 wil help us
// to create multiple buckets using modulo
func (m *Memory) generateHash(keyword []byte) uint64 {
	return crc64.Checksum(keyword, m.crcTable)
}
