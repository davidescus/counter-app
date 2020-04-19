package storage

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	dataFile = "base"
)

type Disk struct {
	path string
}

func NewDisk(path string) (*Disk, error) {
	disk := Disk{
		path: path,
	}
	return &disk, disk.initiateFiles()
}

// Load will load data into memory after each start
func (d *Disk) Load(memory *Memory) error {
	log.Println("--- TODO load")

	path, err := filepath.Abs(d.path)
	if err != nil {
		return err
	}
	f, err := os.Open(path + "/" + dataFile)
	if err != nil {
		return err
	}
	defer f.Close()

	// TODO will refactor this when add algorithm for consistency
	// memory layout: key => value 16bytes with no separator
	var data [16]byte
	for {
		_, err := f.Read(data[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		key, err := d.decode(data[0:8])
		if err != nil {
			return err
		}
		val, err := d.decode(data[8:16])
		if err != nil {
			return err
		}

		// we assume we have only one time a key on disk
		memory.occurrences[key] = val
	}

	return nil
}

func (d *Disk) encode(i uint64) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], i)
	return b[:]
}

func (d *Disk) decode(b []byte) (uint64, error) {
	var i uint64
	reader := bytes.NewReader(b)
	err := binary.Read(reader, binary.LittleEndian, &i)
	return i, err
}

// Flush will create snapshot on disk
func (d *Disk) Flush(m *Memory) error {
	path, err := filepath.Abs(d.path)
	if err != nil {
		return err
	}

	// TODO better permission management, why 0644 not works?
	f, err := os.OpenFile(path+"/"+dataFile, os.O_WRONLY, 0766)
	defer f.Close()
	if err != nil {
		return err
	}

	// TODO apply lock here, copy values, etc improve this
	// TODO make a copy before
	var data []byte
	for k, v := range m.occurrences {
		data = append(data, d.encode(k)...)
		data = append(data, d.encode(v)...)
	}
	_, err = f.Write(data)

	return err
}

func (d *Disk) initiateFiles() error {
	path, err := filepath.Abs(d.path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0766); err != nil {
			return err
		}
	}
	// TODO better permission management, why 0644 not works?
	if _, err := os.OpenFile(path+"/"+dataFile, os.O_APPEND|os.O_CREATE, 0766); err != nil {
		return err
	}

	return nil
}
