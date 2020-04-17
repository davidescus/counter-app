package storage

import "sync"

type Memory struct {
	mu *sync.RWMutex
	occurrences map[string]int64
}

func NewMemory() *Memory {
	return &Memory{
		occurrences: make(map[string]int64),
	}
}

func (m *Memory) Set(key []byte) {
   m.occurrences[string(key)]++
}

func (m *Memory) Get(key []byte) int64 {
	return m.occurrences[string(key)]
}
