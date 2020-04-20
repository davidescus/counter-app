package syncapi

// Storage ...
type Storage interface {
	Merge(occurrences map[uint64][]uint64)
}

type Occurrence struct {
	HashKeyword uint64   `json:"hash_keyword"`
	Totals      []uint64 `json:"totals"`
}
