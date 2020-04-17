package app

import "testing"

// TODO add benchmarks

func TestApp_StoreGet(t *testing.T) {
	table := []struct {
		name   string
		keys   []string
		search []string
		exp []Occurrence
	}{
		{
			name:   "noValuesNoOccurrence",
			keys:   []string{},
			search: []string{},
			exp: []Occurrence{},
		},
		{
			name:   "multipleValuesNoOccurrence",
			keys:   []string{"value1", "value2"},
			search: []string{"one"},
			exp: []Occurrence{
				{Key: "one", Count: uint64(0)},
			},
		},
		{
			name:   "multipleValuesLessOccurrences",
			keys:   []string{"one", "one", "two"},
			search: []string{"one"},
			exp: []Occurrence{
				{Key: "one", Count: uint64(2)},
			},
		},
		{
			name:   "multipleValuesMultipleOccurrences",
			keys:   []string{"one", "one", "two", "four", "two"},
			search: []string{"one", "two"},
			exp: []Occurrence{
				{Key: "one", Count: uint64(2)},
				{Key: "two", Count: uint64(2) },
			},
		},
		{
			name:   "valuesEqualWithOccurrencesNumber",
			keys:   []string{"three", "three", "three"},
			search: []string{"three"},
			exp: []Occurrence{
				{Key: "one", Count: uint64(3)},
			},
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			a := New()
			for _, key := range tt.keys {
				_ = a.Store(key)
			}
			got := a.Get(tt.search)
			if !hasSameOccurrencesNumber(tt.exp, got) {
				t.Errorf("Exp: %v, Got: %v", tt.exp, got)
			}
		})
	}
}

func hasSameOccurrencesNumber(s1, s2 []Occurrence) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
        if s1[i].Count != s2[i].Count {
        	return false
		}
	}
	return true
}