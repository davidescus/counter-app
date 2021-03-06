package memory

import (
	"context"
	"log"
	"os"
	"testing"
)

// TODO add benchmarks
var (
	logger = log.New(os.Stdout, "", 0)
	ctx    = context.Background()
	conf   = &Config{ID: 4}
)

func TestMemory_SetGet(t *testing.T) {
	table := []struct {
		name   string
		keys   []string
		expKey string
		expVal int
	}{
		{
			name:   "noValuesNoOccurrence",
			keys:   []string{},
			expKey: "one",
			expVal: 0,
		},
		{
			name:   "multipleValuesNoOccurrence",
			keys:   []string{"value1", "value2"},
			expKey: "one",
			expVal: 0,
		},
		{
			name:   "multipleValuesLessOccurrences",
			keys:   []string{"one", "one", "two"},
			expKey: "one",
			expVal: 2,
		},
		{
			name:   "valuesEqualWithOccurrencesNumber",
			keys:   []string{"three", "three", "three"},
			expKey: "three",
			expVal: 3,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMemory(ctx, logger, conf)
			for _, key := range tt.keys {
				m.Increment([]byte(key))
			}
			got := m.Get([]byte(tt.expKey))
			if got != uint64(tt.expVal) {
				t.Errorf("Exp: %d, Got: %d", tt.expVal, got)
			}
		})
	}
}
