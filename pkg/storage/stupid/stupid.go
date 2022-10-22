package stupid

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/utkarsh-pro/use/pkg/storage/tlvrw"
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

	rmu *sync.Mutex
	rfd *os.File

	wmu *sync.Mutex
	wfd *os.File
}

// New returns a new Storage instance.
func New(dir string) *Storage {
	return &Storage{
		file: filepath.Join(dir, "stupid.db"),
		rmu:  &sync.Mutex{},
		wmu:  &sync.Mutex{},
	}
}

// Init configures the storage.
func (s *Storage) Init() error {
	wfd, err := os.OpenFile(s.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening file for write fd: %w", err)
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

	s.rmu.Lock()
	defer func() {
		if err := ResetReader(s.rfd); err != nil {
			// fail early
			panic(fmt.Errorf("error resetting reader: %w", err))
		}

		s.rmu.Unlock()
	}()

	var candidate *PacketLite = nil

	for {
		packet, err := ReadPacketLite(s.rfd)
		if err != nil {
			if errors.Is(err, tlvrw.EOF) {
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

	val, err := readValueAt(s.rfd, candidate.ValPos)
	if err != nil {
		return nil, fmt.Errorf("error reading value: %w", err)
	}

	return val, nil
}

// Set sets the value for the given key.
func (s *Storage) Set(key string, value []byte) error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	if err := WritePacket(
		s.wfd,
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

	return nil
}

// Delete deletes the value for the given key.
func (s *Storage) Delete(key string) error {
	if !s.isInit() {
		return ErrStorageNotInitialized
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	if err := WritePacket(
		s.wfd,
		&Packet{
			Op:  DelOp,
			Key: []byte(key),
			Val: nil,
		},
	); err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}

	return nil
}

// Exists returns true if the given key exists.
func (s *Storage) Exists(key string) (bool, error) {
	if !s.isInit() {
		return false, ErrStorageNotInitialized
	}

	s.rmu.Lock()
	defer func() {
		if err := ResetReader(s.rfd); err != nil {
			// fail early
			panic(fmt.Errorf("error resetting reader: %w", err))
		}

		s.rmu.Unlock()
	}()

	var candidate *PacketLite = nil

	for {
		packet, err := ReadPacketLite(s.rfd)
		if err != nil {
			if errors.Is(err, tlvrw.EOF) {
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

	s.rmu.Lock()
	defer s.rmu.Unlock()

	if err := s.rfd.Close(); err != nil {
		return fmt.Errorf("error closing read fd: %w", err)
	}
	s.rfd = nil

	s.wmu.Lock()
	defer s.wmu.Unlock()

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

	s.rmu.Lock()
	defer func() {
		if err := ResetReader(s.rfd); err != nil {
			// fail early
			panic(fmt.Errorf("error resetting reader: %w", err))
		}

		s.rmu.Unlock()
	}()

	set := make(map[string]struct{})

	for {
		packet, err := ReadPacketLite(s.rfd)
		if err != nil {
			if errors.Is(err, tlvrw.EOF) {
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

// isInit returns true if the storage is initialized.
func (s *Storage) isInit() bool {
	return s.rfd != nil && s.wfd != nil
}
