package storage

import (
	"testing"
)

// TODO add tests for flush and load, use /tmp as storagePath

func TestDisk_EncodeDecode(t *testing.T) {
	table := []struct {
		name string
		exp  uint64
	}{
		{
			name: "zeroValue",
			exp:  uint64(0),
		},
		{
			name: "smallNumber",
			exp:  uint64(99),
		},
		{
			name: "bigNumber",
			exp:  uint64(9999999999),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			d := Disk{}
			got, _ := d.decode(d.encode(tt.exp))
			if got != tt.exp {
				t.Errorf("Exp: %d, Got: %d", tt.exp, got)
			}
		})
	}
}
