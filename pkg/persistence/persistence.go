package persistence

// Storage represent contract with in memory storage,
// persistent storage need only this in fro in memory API
type Storage interface {
	Import(map[uint64][]uint64)
	Export() map[uint64][]uint64
}

// PersistentStorage represents API expose outside
type PersistentStorage interface {
	LoadVolatileStorage() error
	FlushNow() error
	StartFlashing()
	StopFlashing() error
}
