package localdisk

import (
	"context"
	"counter-app/pkg/persistence"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Conf struct {
	Path, File      string
	FlushIntervalMS int
}

type LocalDisk struct {
	ctx             context.Context
	logger          *log.Logger
	volatileStorage persistence.Storage
	path            string
	file            string
	flashIntervalMS int
	finishFlush     chan struct{}
}

func NewLocalDisk(ctx context.Context, logger *log.Logger, conf *Conf, memory persistence.Storage) (*LocalDisk, error) {
	ld := LocalDisk{
		ctx:             ctx,
		logger:          logger,
		volatileStorage: memory,
		path:            conf.Path,
		file:            conf.File,
		flashIntervalMS: conf.FlushIntervalMS,
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
	// TODO implement me
	var data []byte
	d.volatileStorage.Import(data)

	return nil
}

// FlushNow will flush to disk all data now
func (d *LocalDisk) FlushNow() error {
	log.Println("--- FlushNow()")
	data := d.volatileStorage.Export()

	// TODO implement logic here
	_ = data

	return nil
}

// StartFlashing will start a goroutine that will
// flush data to disc periodic at specific interval
func (d *LocalDisk) StartFlashing() {
	go func(ch chan struct{}) {
		for {
			select {
			case <-d.ctx.Done():
				ch <- struct{}{}
				return
			default:
				if err := d.FlushNow(); err != nil {
					d.logger.Println("[Error] Flash to dish: ", err)
				}
				time.Sleep(time.Duration(d.flashIntervalMS) * time.Millisecond)
			}
		}
	}(d.finishFlush)
}

// StopFlashing will wait for finish flush signal
// and execute last flash on disk
func (d *LocalDisk) StopFlashing() error {
	log.Println("-----Wait for finish flush")
	<-d.finishFlush
	log.Println("------ finish stop flush")
	return d.FlushNow()

}
