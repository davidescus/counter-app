package localdisk

import (
	"bytes"
	"context"
	"counter-app/pkg/persistence"
	"encoding/binary"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Conf struct {
	Path, File   string
	DumpInterval time.Duration
}

type LocalDisk struct {
	ctx                   context.Context
	logger                *log.Logger
	volatileStorage       persistence.Storage
	path                  string
	file                  string
	dumpInterval          time.Duration
	endSig, endSigSuccess chan struct{}
}

func NewLocalDisk(ctx context.Context, logger *log.Logger, conf *Conf, memory persistence.Storage) (*LocalDisk, error) {
	ld := LocalDisk{
		ctx:             ctx,
		logger:          logger,
		volatileStorage: memory,
		path:            conf.Path,
		file:            conf.File,
		dumpInterval:    conf.DumpInterval,
		endSig:          make(chan struct{}),
		endSigSuccess:   make(chan struct{}),
	}
	return &ld, ld.initiateFiles()
}

func (d *LocalDisk) initiateFiles() error {
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
	if _, err := os.OpenFile(path+"/"+d.file, os.O_APPEND|os.O_CREATE, 0766); err != nil {
		return err
	}

	return nil
}

func (d *LocalDisk) LoadVolatileStorage() error {
	data, err := d.ReadFromDisk()
	if err != nil {
		return err
	}
	d.volatileStorage.Import(data)

	return nil
}

// DumpNow will flush to disk all data now
func (d *LocalDisk) DumpNow() error {
	return d.WriteToDisk(d.volatileStorage.Export())
}

// StartPersistentDumps will start a goroutine that will
// flush data to disc periodic at specific interval
func (d *LocalDisk) StartPersistentDumps() {
	go func() {
		for {
			select {
			case <-d.ctx.Done():
			case <-d.endSig:
				d.endSigSuccess <- struct{}{}
				return
			default:
				if err := d.DumpNow(); err != nil {
					d.logger.Println("[Error] PersistenceDump: ", err)
				}
				time.Sleep(d.dumpInterval)
			}
		}
	}()
}

// StopPersistenceDump will wait for finish flush signal
// and execute last Dump on disk
func (d *LocalDisk) StopPersistenceDump() error {
	d.endSig <- struct{}{}
	<-d.endSigSuccess
	return d.DumpNow()
}

// WriteToDisk will write data on disk
func (d *LocalDisk) WriteToDisk(data map[uint64][]uint64) error {
	path, err := filepath.Abs(d.path)
	if err != nil {
		return err
	}

	// TODO better permission management, why 0644 not works?
	f, err := os.OpenFile(path+"/"+d.file, os.O_WRONLY, 0766)
	defer f.Close()
	if err != nil {
		return err
	}

	// TODO create new file before write, swap it after

	var buf []byte
	var placeholder [8]byte
	// memory layout: uint64(key) uint64(len(totals) [uint64 ...](totals values)
	for k, v := range data {
		binary.LittleEndian.PutUint64(placeholder[:], k)
		buf = append(buf, placeholder[:]...)
		// waste resources with uint64, but keep it simple at the moment
		binary.LittleEndian.PutUint64(placeholder[:], uint64(len(v)))
		buf = append(buf, placeholder[:]...)

		// v []uint64, range over it and add them
		for _, vv := range v {
			binary.LittleEndian.PutUint64(placeholder[:], vv)
			buf = append(buf, placeholder[:]...)
		}
	}

	_, err = f.Write(buf)

	return err
}

// ReadFromDisk will read data and map it into struncture
func (d *LocalDisk) ReadFromDisk() (map[uint64][]uint64, error) {
	path, err := filepath.Abs(d.path)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path + "/" + d.file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data := make(map[uint64][]uint64)
	var placeholder [8]byte
	// memory layout: uint64(key) uint64(len(totals) [uint64 ...](totals values)
	for {
		// hash
		_, err := f.Read(placeholder[:])
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		key, err := decode(placeholder[0:8])
		if err != nil {
			return nil, err
		}

		// keys len
		if _, err := f.Read(placeholder[:]); err != nil {
			return nil, err
		}
		sliceLen, err := decode(placeholder[0:8])
		if err != nil {
			return nil, err
		}

		// get totals
		var i uint64
		var totals []uint64
		for ; i < sliceLen; i++ {
			if _, err := f.Read(placeholder[:]); err != nil {
				return nil, err
			}
			v, err := decode(placeholder[0:8])
			if err != nil {
				return nil, err
			}
			totals = append(totals, v)
		}

		data[key] = totals
	}

	return data, nil
}

func decode(b []byte) (uint64, error) {
	var i uint64
	reader := bytes.NewReader(b)
	err := binary.Read(reader, binary.LittleEndian, &i)
	return i, err
}
