package config

// SyncType is the type of sync.
type SyncType string

var (
	// SyncTypeNone is the none sync type.
	SyncTypeNone SyncType = "none"
	// SyncTypeSync is the sync sync type.
	SyncTypeSync SyncType = "sync"
	// SyncTypeAsync is the async sync type.
	SyncTypeAsync SyncType = "async"
)

// Config is the config for the storage.
type Config struct {
	Sync     SyncType
	ReadOnly bool
}

// DefaultConfig returns the default config.
func DefaultConfig() Config {
	return Config{
		Sync: SyncTypeNone,
	}
}

// WithSync sets the sync type.
func (cfg Config) WithSync() Config {
	cfg.Sync = SyncTypeSync
	return cfg
}

// WithAsyncSync sets the async type.
func (cfg Config) WithAsyncSync() Config {
	cfg.Sync = SyncTypeAsync
	return cfg
}

// WithNoneSync sets the none type.
func (cfg Config) WithNoneSync() Config {
	cfg.Sync = SyncTypeNone
	return cfg
}

// WithReadOnly sets the read only type.
func (cfg Config) WithReadOnly() Config {
	cfg.ReadOnly = true
	return cfg
}

// WithReadWrite sets the read write type.
func (cfg Config) WithReadWrite() Config {
	cfg.ReadOnly = false
	return cfg
}
