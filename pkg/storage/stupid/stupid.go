package stupid

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var (
	// Storage Ops
	SetOp = byte(1)
	DelOp = byte(2)

	// Errors
	ErrStorageNotInitialized = fmt.Errorf("storage is not initialized")
	ErrKeyNotFound           = fmt.Errorf("key not found")
)

// Storage is a stupid storage.
type Storage struct {
	file string

	rfd *os.File

	wmu                 *sync.Mutex
	wfd                 *os.File
	lastSuccessWritePos int64
}

// New returns a new Storage instance.
func New(dir string) *Storage {
	return &Storage{
		file: filepath.Join(dir, "stupid.db"),
		wmu:  &sync.Mutex{},
	}
}

// Init configures the storage.
func (s *Storage) Init() error {
	wfd, err := os.OpenFile(s.file, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening file for write fd: %w", err)
	}
	// Move to the end of the file.
	if pos, err := wfd.Seek(0, io.SeekEnd); err != nil {
		return fmt.Errorf("error seeking to end of file: %w", err)
	} else {
		s.lastSuccessWritePos = pos
	}

	rfd, err := os.OpenFile(s.file, os.O_RDONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening file for read fd: %w", err)
	}

	s.rfd = rfd
	s.wfd = wfd

	return nil
}

// Get returns the value for the given key.
func (s *Storage) Get(key string) ([]byte, error) {
	if !s.isInit() {
		return nil, ErrStorageNotInitialized
	}

	var candidate *Packet = nil
	pr := newreader(s.rfd)

	for {
		packet := &Packet{}
		if err := pr.lread(packet); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("error reading packet: %w", err)
		}

		if string(packet.Key) == key {
			if packet.Op == DelOp {
				candidate = nil
			}

			if packet.Op == SetOp {
				candidate = packet
			}
		}
	}

	if candidate == nil {
		return nil, ErrKeyNotFound
	}

	if err := pr.fill(candidate); err != nil {
		return nil, fmt.Errorf("error filling packet: %w", err)
	}

	return candidate.Val, nil
}

// Set sets the value for the given key.
func (s *Storage) Set(key string, value []byte) error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	pw := newwriter(s.wfd)
	if err := pw.write(
		&Packet{
			Op:  SetOp,
			Key: []byte(key),
			Val: value,
		},
	); err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}

	if err := s.wfd.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %w", err)
	}

	// record the last successful write position
	pos, err := s.wfd.Seek(0, io.SeekCurrent)
	if err != nil {
		fmt.Println("[WARN]: failed to get current write position: ", err)
		return nil
	}

	s.lastSuccessWritePos = pos
	return nil
}

// Delete deletes the value for the given key.
func (s *Storage) Delete(key string) error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	pw := newwriter(s.wfd)
	if err := pw.write(
		&Packet{
			Op:  DelOp,
			Key: []byte(key),
			Val: nil,
		},
	); err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}

	if err := s.wfd.Sync(); err != nil {
		return fmt.Errorf("error syncing file: %w", err)
	}

	// record the last successful write position
	pos, err := s.wfd.Seek(0, io.SeekCurrent)
	if err != nil {
		fmt.Println("[WARN]: failed to get current write position: ", err)
		return nil
	}

	s.lastSuccessWritePos = pos
	return nil
}

// Exists returns true if the given key exists.
func (s *Storage) Exists(key string) (bool, error) {
	if !s.isInit() {
		return false, ErrStorageNotInitialized
	}

	var candidate *Packet = nil
	pr := newreader(s.rfd)

	for {
		packet := &Packet{}
		if err := pr.lread(packet); err != nil {
			if err == io.EOF {
				break
			}

			return false, fmt.Errorf("error reading packet: %w", err)
		}

		if string(packet.Key) == key {
			if packet.Op == DelOp {
				candidate = nil
			}

			if packet.Op == SetOp {
				candidate = packet
			}
		}
	}

	return candidate != nil, nil
}

// Close closes the storage.
func (s *Storage) Close() error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	if err := s.rfd.Close(); err != nil {
		return fmt.Errorf("error closing read fd: %w", err)
	}
	s.rfd = nil

	if err := s.wfd.Close(); err != nil {
		return fmt.Errorf("error closing write fd: %w", err)
	}
	s.wfd = nil

	return nil
}

func (s *Storage) Len() (int, error) {
	if !s.isInit() {
		return 0, ErrStorageNotInitialized
	}

	set := make(map[string]struct{})
	pr := newreader(s.rfd)

	for {
		packet := &Packet{}
		if err := pr.lread(packet); err != nil {
			if err == io.EOF {
				break
			}

			return 0, fmt.Errorf("error reading packet: %w", err)
		}

		if packet.Op == DelOp {
			delete(set, string(packet.Key))
			continue
		}

		if packet.Op == SetOp {
			set[string(packet.Key)] = struct{}{}
			continue
		}
	}

	return len(set), nil
}

// PhysicalSnapshot writes the current state of the storage to the given writer.
// The snapshot is written in a format that can be read by only by the storage.
//
// Note: The DB is locked for writes while the snapshot is being generated.
func (s *Storage) PhysicalSnapshot(w io.Writer) error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	// Get the last successful write position
	lastSuccessWritePos := s.lastSuccessWritePos

	// Create a new reader for the file.
	tempfd, err := os.Open(s.file)
	if err != nil {
		return fmt.Errorf("error opening file for read fd: %w", err)
	}

	if _, err := io.CopyN(w, tempfd, lastSuccessWritePos); err != nil {
		return fmt.Errorf("error generating snapshot: %w", err)
	}

	return nil
}

// isInit returns true if the storage is initialized.
func (s *Storage) isInit() bool {
	return s.rfd != nil && s.wfd != nil
}
