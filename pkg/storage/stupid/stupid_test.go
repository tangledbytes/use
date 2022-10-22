package stupid

import (
	"testing"
)

func TestAll(t *testing.T) {
	dir := t.TempDir()

	s := New(dir)

	t.Run("isInit", func(t *testing.T) {
		if s.isInit() {
			t.Error("storage is initialized")
		}

		t.Run("Get", func(t *testing.T) {
			_, err := s.Get("foo")
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Set", func(t *testing.T) {
			err := s.Set("foo", nil)
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := s.Delete("foo")
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Exists", func(t *testing.T) {
			_, err := s.Exists("foo")
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Len", func(t *testing.T) {
			_, err := s.Len()
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Close", func(t *testing.T) {
			err := s.Close()
			if err != ErrStorageNotInitialized {
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

	t.Run("Set", func(t *testing.T) {
		t.Run("Valid Set", func(t *testing.T) {
			if err := s.Set("foo", []byte("bar")); err != nil {
				t.Error(err)
			}

			if err := s.Set("foo", []byte("bazz")); err != nil {
				t.Error(err)
			}

			if err := s.Set("foo2", []byte("bazz2")); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("Valid Get", func(t *testing.T) {
			val, err := s.Get("foo")
			if err != nil {
				t.Error(err)
			}

			if string(val) != "bazz" {
				t.Error("expected bazz")
			}

			val, err = s.Get("foo2")
			if err != nil {
				t.Error(err)
			}

			if string(val) != "bazz2" {
				t.Error("expected bazz2")
			}
		})

		t.Run("Invalid Get", func(t *testing.T) {
			_, err := s.Get("foo3")
			if err != ErrKeyNotFound {
				t.Error("expected ErrKeyNotFound")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("Valid Delete", func(t *testing.T) {
			if err := s.Delete("foo"); err != nil {
				t.Error(err)
			}

			if _, err := s.Get("foo"); err != ErrKeyNotFound {
				t.Error("expected ErrKeyNotFound")
			}
		})
	})

	t.Run("Exists", func(t *testing.T) {
		t.Run("Valid Exists", func(t *testing.T) {
			if ok, err := s.Exists("foo2"); err != nil {
				t.Error(err)
			} else if !ok {
				t.Error("expected true")
			}

			if ok, err := s.Exists("foo3"); err != nil {
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
			} else if n != 1 {
				t.Error("expected 1")
			}

			if err := s.Set("foo3", []byte("bazz")); err != nil {
				t.Error(err)
			}

			if n, err := s.Len(); err != nil {
				t.Error(err)
			} else if n != 2 {
				t.Error("expected 2")
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
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Set", func(t *testing.T) {
			err := s.Set("foo", nil)
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Delete", func(t *testing.T) {
			err := s.Delete("foo")
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}

		})

		t.Run("Exists", func(t *testing.T) {
			_, err := s.Exists("foo")
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})

		t.Run("Len", func(t *testing.T) {
			_, err := s.Len()
			if err != ErrStorageNotInitialized {
				t.Error("expected ErrStorageNotInitialized")
			}
		})
	})
}
