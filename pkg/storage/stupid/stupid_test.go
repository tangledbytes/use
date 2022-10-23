package stupid

import (
	"math/rand"
	"testing"
	"time"

	"github.com/utkarsh-pro/use/pkg/storage/errors"
)

func generateRandomBytes(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	rand.Read(b)

	return b
}

func TestAll(t *testing.T) {
	dir := t.TempDir()

	s := New(dir)

	t.Run("isInit", func(t *testing.T) {
		if s.isInit() {
			t.Error("storage is initialized")
		}

		t.Run("Get", func(t *testing.T) {
			_, err := s.Get("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Set", func(t *testing.T) {
			err := s.Set("foo", nil)
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := s.Delete("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Exists", func(t *testing.T) {
			_, err := s.Exists("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Len", func(t *testing.T) {
			_, err := s.Len()
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Close", func(t *testing.T) {
			err := s.Close()
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})
	})

	t.Run("Init", func(t *testing.T) {
		if err := s.Init(); err != nil {
			t.Error(err)
		}

		t.Run("isInit", func(t *testing.T) {
			if !s.isInit() {
				t.Error("storage is not initialized")
			}
		})
	})

	valsforSet := make(map[string][]byte)
	valsforSet["foo"] = []byte("bar")
	valsforSet["bar"] = []byte("baz")
	valsforSet["baz"] = []byte("foo")
	valsforSet["mr.big.empty"] = make([]byte, 1024*1024*10)
	valsforSet["mr.big.random"] = generateRandomBytes(1024 * 1024 * 10)

	t.Run("Set", func(t *testing.T) {
		t.Run("Valid Set", func(t *testing.T) {
			for k, v := range valsforSet {
				if err := s.Set(k, v); err != nil {
					t.Error(err)
				}
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Valid Get", func(t *testing.T) {
			for k, v := range valsforSet {
				val, err := s.Get(k)
				if err != nil {
					t.Error(err)
				}

				if string(val) != string(v) {
					t.Error("value mismatch")
				}
			}
		})

		t.Run("Invalid Get", func(t *testing.T) {
			_, err := s.Get("foo3")
			if err != errors.ErrKeyNotFound {
				t.Error("expected ErrKeyNotFound")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Valid Delete", func(t *testing.T) {
			// Delete the first key that we encounter
			for k := range valsforSet {
				if err := s.Delete(k); err != nil {
					t.Error(err)
				}

				_, err := s.Get(k)
				if err != errors.ErrKeyNotFound {
					t.Error("expected ErrKeyNotFound")
				}

				delete(valsforSet, k)
				return
			}
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("Valid Exists", func(t *testing.T) {
			for k := range valsforSet {
				exists, err := s.Exists(k)
				if err != nil {
					t.Error(err)
				}

				if !exists {
					t.Error("key does not exist")
				}
			}

			if ok, err := s.Exists(string(generateRandomBytes(5))); err != nil {
				t.Error(err)
			} else if ok {
				t.Error("expected false")
			}
		})
	})

	t.Run("Len", func(t *testing.T) {
		t.Run("Valid Len", func(t *testing.T) {
			if n, err := s.Len(); err != nil {
				t.Error(err)
			} else if n != len(valsforSet) {
				t.Error("expected", len(valsforSet), "got", n)
			}

			if err := s.Set("foo3", []byte("bazz")); err != nil {
				t.Error(err)
			}

			if n, err := s.Len(); err != nil {
				t.Error(err)
			} else if n != len(valsforSet)+1 {
				t.Error("expected", len(valsforSet)+1, "got", n)
			}
		})
	})

	t.Run("Close", func(t *testing.T) {
		if err := s.Close(); err != nil {
			t.Error(err)
		}

		t.Run("isInit", func(t *testing.T) {
			if s.isInit() {
				t.Error("storage is initialized")
			}
		})

		t.Run("Get", func(t *testing.T) {
			_, err := s.Get("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Set", func(t *testing.T) {
			err := s.Set("foo", nil)
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := s.Delete("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}

		})

		t.Run("Exists", func(t *testing.T) {
			_, err := s.Exists("foo")
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Len", func(t *testing.T) {
			_, err := s.Len()
			if err != errors.ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})
	})
}
