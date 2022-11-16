package stupid

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/utkarsh-pro/use/pkg/id"
	"github.com/utkarsh-pro/use/pkg/log"
	"github.com/utkarsh-pro/use/pkg/storage/config"
	"github.com/utkarsh-pro/use/pkg/storage/errors"
	"github.com/utkarsh-pro/use/pkg/structures/bloom/dibf"
)

var (
	// Storage Ops
	SetOp = byte(1)
	DelOp = byte(2)
)

// Storage is a stupid storage.
type Storage struct {
	// file is path to the storage file.
	file string

	// rfd is the reading file descriptor.
	rfd *os.File

	// idgen is the id generator.
	idgen id.Gen

	// wmu is the write mutex.
	wmu *sync.Mutex
	// wfd is the writing file descriptor.
	wfd *os.File
	// lastSuccessWritePos is the last successful write position.
	lastSuccessWritePos int64

	// cfg is the storage config.
	cfg config.Config

	// bf is the bloom bf.
	bf bf
}

// New returns a new Storage instance.
func New(dir string, cfg config.Config) *Storage {
	return &Storage{
		file:  filepath.Join(dir, "stupid.db"),
		wmu:   &sync.Mutex{},
		idgen: id.New(),
		cfg:   cfg,
		bf:    newBfSync(dibf.NewWithEstimates(1e6, 0.01, 1, nil)),
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

	// Don't perform any recovery if the storage is read-only.
	if s.cfg.ReadOnly {
		log.Infoln("storage is read-only, skipping recovery")
		return nil
	}

	// Fix the corrupt data if there is any
	if err := s.DetectAndFix(); err != nil {
		return fmt.Errorf("%s: %w", errors.ErrCorruptStorage, err)
	}

	return nil
}

// Get returns the value for the given key.
func (s *Storage) Get(key string) ([]byte, error) {
	if !s.isInit() {
		return nil, errors.ErrStorageNotInitialized
	}

	if !s.bf.Contains([]byte(key)) {
		return nil, errors.ErrKeyNotFound
	}

	var candidate *Packet = nil
	var pr *reader = nil

	err := s.ForEach(func(r *reader, p *Packet, err error) error {
		// get the reader
		if pr == nil {
			pr = r
		}

		if err != nil {
			// Although this should NEVER happen. NEVER.
			if err == io.EOF {
				return nil
			}
		}

		if string(p.Key) == key {
			if p.Op == DelOp {
				candidate = nil
				return nil
			}

			if p.Op == SetOp {
				candidate = p
				return nil
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if candidate == nil {
		return nil, errors.ErrKeyNotFound
	}

	if err := pr.fill(candidate); err != nil {
		return nil, err
	}

	return candidate.Val, nil
}

// Set sets the value for the given key.
func (s *Storage) Set(key string, value []byte) error {
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
	}

	if s.cfg.ReadOnly {
		return errors.ErrReadOnlyStorage
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	pw := newwriter(s.wfd)
	if err := pw.write(
		&Packet{
			ID:  s.idgen.Next(),
			Op:  SetOp,
			Key: []byte(key),
			Val: value,
		},
	); err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}

	if s.cfg.Sync == config.SyncTypeSync {
		if err := s.wfd.Sync(); err != nil {
			return fmt.Errorf("error syncing file: %w", err)
		}
	} else if s.cfg.Sync == config.SyncTypeAsync {
		go s.wfd.Sync()
	}

	// record the last successful write position
	pos, err := s.wfd.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Warnln("failed to get current write position: ", err)
		return nil
	}

	s.lastSuccessWritePos = pos

	// add to bloom filter
	s.bf.Add([]byte(key))

	return nil
}

// Delete deletes the value for the given key.
func (s *Storage) Delete(key string) error {
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
	}

	if s.cfg.ReadOnly {
		return errors.ErrReadOnlyStorage
	}

	if !s.bf.Contains([]byte(key)) {
		return nil
	}

	s.wmu.Lock()
	defer s.wmu.Unlock()

	pw := newwriter(s.wfd)
	if err := pw.write(
		&Packet{
			ID:  s.idgen.Next(),
			Op:  DelOp,
			Key: []byte(key),
			Val: nil,
		},
	); err != nil {
		return fmt.Errorf("error writing packet: %w", err)
	}

	if s.cfg.Sync == config.SyncTypeSync {
		if err := s.wfd.Sync(); err != nil {
			return fmt.Errorf("error syncing file: %w", err)
		}
	} else if s.cfg.Sync == config.SyncTypeAsync {
		go s.wfd.Sync()
	}

	// record the last successful write position
	pos, err := s.wfd.Seek(0, io.SeekCurrent)
	if err != nil {
		log.Warnln("failed to get current write position: ", err)
		return nil
	}

	s.lastSuccessWritePos = pos

	// remove from bloom filter
	s.bf.Delete([]byte(key))
	return nil
}

// Exists returns true if the given key exists.
func (s *Storage) Exists(key string) (bool, error) {
	if !s.isInit() {
		return false, errors.ErrStorageNotInitialized
	}

	return s.bf.Contains([]byte(key)), nil
}

// Len returns the number of keys in the storage.
func (s *Storage) Len() (int, error) {
	if !s.isInit() {
		return 0, errors.ErrStorageNotInitialized
	}

	set := make(map[string]struct{})

	if err := s.ForEach(func(r *reader, p *Packet, err error) error {
		if err != nil {
			return err
		}

		set[string(p.Key)] = struct{}{}
		return nil
	}); err != nil {
		return 0, err
	}

	return len(set), nil
}

// PhysicalSnapshot writes the current state of the storage to the given writer.
// The snapshot is written in a format that can be read by only by the storage.
//
// Note: The DB is locked for writes while the snapshot is being generated.
func (s *Storage) PhysicalSnapshot(w io.Writer) error {
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
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

// ForEach goes through the entire store and executes the given function
// on each packet that it reads.
func (s *Storage) ForEach(fn func(*reader, *Packet, error) error) error {
	// Fail silently if the storage hasn't been initialized yet
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
	}

	pr := newreader(s.rfd)

	for {
		// don't read beyond the last successful write position
		if pr.pos() >= s.lastSuccessWritePos {
			break
		}

		packet := &Packet{}
		err := pr.lread(packet)
		if err != nil {
			// If we reach the end of the file, break.
			//
			// We don't call the function with the packet
			// because we know that the packet has to be invalid.
			//
			// Due to the design of TLV, it is impossible to encounter
			// EOF while reading a packet (even at the end).
			if err == io.EOF {
				break
			}
		}

		if _err := fn(pr, packet, err); _err != nil {
			return _err
		}
	}

	return nil
}

// DetectAndFix detects corrupt data in the store and tries to fix it
//
// Note: DetectAndFix will remove the corrupt data from the store which means
// that some of the writes might vanish. But this is same because the storage engine
// provides guaranatees only in the cases when it is running in the sync mode and a
// response is sent back to the client indicating a success.
func (s *Storage) DetectAndFix() error {
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
	}

	lastSuccessRead := int64(0)

	// Go through the entire store and try to see if we can get valid
	// packet reads from the store
	return s.ForEach(func(pr *reader, p *Packet, err error) error {
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				log.Warnln("Found corrupted data in the store. Trying to fix it...")

				// Discard the rest of the file
				if err := s.wfd.Truncate(lastSuccessRead); err != nil {
					return fmt.Errorf("error truncating file: %w", err)
				}

				// No need to update the read position since we only rely on ReadAt
				// and that doesn't changes os.File read position

				// Move the write position to the last successful read position
				s.wfd.Seek(lastSuccessRead, io.SeekStart)

				// Update the last successful write position
				s.lastSuccessWritePos = lastSuccessRead

				log.Infoln("Successfully fixed the corrupted data in the store")
				return nil
			}

			return err
		}

		// Update the last successful read position
		lastSuccessRead = pr.pos()

		// Insert the packet into the bloom filter
		s.bf.Add(p.Key)

		return nil
	})
}

// Close closes the storage.
func (s *Storage) Close() error {
	if !s.isInit() {
		return errors.ErrStorageNotInitialized
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

// GetByID returns a packet corresponding to the given ID.
//
// This is a low level API and should not be used by the user.
func (s *Storage) GetByID(id uint64) (*Packet, error) {
	if !s.isInit() {
		return nil, errors.ErrStorageNotInitialized
	}

	var packet *Packet
	desiredErr := fmt.Errorf("desired")

	if err := s.ForEach(func(r *reader, p *Packet, err error) error {
		if err != nil {
			return err
		}

		if p.ID == id {
			packet = p
			return desiredErr
		}

		return nil
	}); err != nil {
		if err == desiredErr {
			return packet, nil
		}

		return nil, err
	}

	return nil, errors.ErrKeyNotFound
}

// isInit returns true if the storage is initialized.
func (s *Storage) isInit() bool {
	return s.rfd != nil && s.wfd != nil
}
