package stupid

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/utkarsh-pro/use/pkg/storage/config"
	"github.com/utkarsh-pro/use/pkg/storage/errors"
	"github.com/utkarsh-pro/use/pkg/utils"
)

func TestAll(t *testing.T) {
	cgfs := map[string]config.Config{
		"nonesynccfg":  config.DefaultConfig(),
		"synccfg":      config.DefaultConfig().WithSync(),
		"asyncsynccfg": config.DefaultConfig().WithAsyncSync(),
	}

	for name, cgf := range cgfs {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()

			s := New(dir, cgf)

			t.Run("isInit", func(t *testing.T) {
				if s.isInit() {
					t.Error("storage is initialized")
				}

				t.Run("Get", func(t *testing.T) {
					_, err := s.Get([]byte("foo"))
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}
				})

				t.Run("Set", func(t *testing.T) {
					err := s.Set([]byte("foo"), nil)
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}
				})

				t.Run("Delete", func(t *testing.T) {
					err := s.Delete([]byte("foo"))
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}
				})

				t.Run("Exists", func(t *testing.T) {
					_, err := s.Exists([]byte("foo"))
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
			valsforSet["mr.big.random"] = utils.GenerateRandomBytes(1024 * 1024 * 10)

			setLength := len(valsforSet)

			t.Run("Set", func(t *testing.T) {
				t.Run("Valid Set", func(t *testing.T) {
					for k, v := range valsforSet {
						if err := s.Set([]byte(k), v); err != nil {
							t.Error(err)
						}
					}
				})
			})

			t.Run("Get", func(t *testing.T) {
				t.Run("Valid Get", func(t *testing.T) {
					for k, v := range valsforSet {
						val, err := s.Get([]byte(k))
						if err != nil {
							t.Error(err)
						}

						if string(val) != string(v) {
							t.Error("value mismatch")
						}
					}
				})

				t.Run("Invalid Get", func(t *testing.T) {
					_, err := s.Get([]byte("foo3"))
					if err != errors.ErrKeyNotFound {
						t.Error("expected ErrKeyNotFound", "got", err)
					}
				})
			})

			t.Run("Delete", func(t *testing.T) {
				t.Run("Valid Delete", func(t *testing.T) {
					// Delete the first key that we encounter
					for k := range valsforSet {
						if err := s.Delete([]byte(k)); err != nil {
							t.Error(err)
						}

						_, err := s.Get([]byte(k))
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
						exists, err := s.Exists([]byte(k))
						if err != nil {
							t.Error(err)
						}

						if !exists {
							t.Error("key does not exist")
						}
					}

					if ok, err := s.Exists(utils.GenerateRandomBytes(5)); err != nil {
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
					} else if n != setLength {
						t.Error("expected", setLength, "got", n)
					}

					if err := s.Set([]byte("foo3"), []byte("bazz")); err != nil {
						t.Error(err)
					}

					if n, err := s.Len(); err != nil {
						t.Error(err)
					} else if n != setLength+1 {
						t.Error("expected", setLength+1, "got", n)
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
					_, err := s.Get([]byte("foo"))
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}
				})

				t.Run("Set", func(t *testing.T) {
					err := s.Set([]byte("foo"), nil)
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}
				})

				t.Run("Delete", func(t *testing.T) {
					err := s.Delete([]byte("foo"))
					if err != errors.ErrStorageNotInitialized {
						t.Error("expected ErrStorageNotInitialized")
					}

				})

				t.Run("Exists", func(t *testing.T) {
					_, err := s.Exists([]byte("foo"))
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
		})
	}
}

func TestStorage_DetectAndFix(t *testing.T) {
	// DetectAndFix test will first create a storage and will write huge chunk of data
	// to it and then will close the storage. After that it will corrupt the storage
	// by removing the last byte of the file and will try to detect and fix the storage.
	// If the storage is fixed successfully, then the test will pass.

	// dir is the directory where the storage will be created.
	dir := t.TempDir()

	// Create a storage
	s := New(dir, config.DefaultConfig())
	if err := s.Init(); err != nil {
		t.Fatal(err)
	}

	// Write huge chunk of data to the storage
	for i := 0; i < 1e5; i++ {
		// A packet is going to be made up of 4 TLVs
		// 1st TLV size => 1 + 4 + 8 = 13
		// 2nd TLV size => 1 + 4 + 1 = 6
		// 3rd TLV size => 1 + 4 + 5 = 10
		// 4th TLV size => 1 + 4 + 5 = 10
		// Total size => 13 + 6 + 10 + 10 = 39
		if err := s.Set(utils.GenerateRandomBytes(5), utils.GenerateRandomBytes(5)); err != nil {
			t.Fatal(err)
		}
	}

	// Close the storage
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name         string
		corruptBytes int
	}{
		{
			name:         "corrupt last byte",
			corruptBytes: 1,
		},
		{
			name:         "corrupt last 2 bytes",
			corruptBytes: 2,
		},
		{
			name:         "corrupt random bytes [1, 40)",
			corruptBytes: utils.GetRandomInRange(1, 40),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Corrupt the storage by removing the last byte of the file
			f, err := os.OpenFile(s.file, os.O_RDWR, 0644)
			if err != nil {
				t.Fatal(err)
			}

			// Get position of the last byte
			fi, err := f.Stat()
			if err != nil {
				t.Fatal(err)
			}

			// Remove the random number of bytes from the end of the file
			if err := f.Truncate(fi.Size() - int64(utils.GetRandomInRange(1, 40))); err != nil {
				t.Fatal(err)
			}

			// Close the file
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}

			// create the store again
			s = New(dir, config.DefaultConfig())

			// Detect and fix the storage
			if err := s.Init(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

type benchmarkTestCase struct {
	name string
	size int
}

func BenchmarkStupidSetNoSync(b *testing.B) {
	dir := b.TempDir()
	s := New(dir, config.DefaultConfig())
	if err := s.Init(); err != nil {
		b.Error(err)
	}
	defer s.Close()

	tests := []benchmarkTestCase{
		{"128B", 128},
		{"256B", 256},
		{"1K", 1024},
		{"2K", 2048},
		{"4K", 4096},
		{"8K", 8192},
		{"16K", 16384},
		{"32K", 32768},
		{"64K", 65536},
		{"128K", 131072},
		{"256K", 262144},
		{"512K", 524288},
		{"1M", 1048576},
		{"2M", 2097152},
		{"4M", 4194304},
		{"8M", 8388608},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.SetBytes(int64(test.size))

			key := "foo"
			value := []byte(strings.Repeat(" ", test.size))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := s.Set([]byte(key), value); err != nil {
					b.Error(err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkStupidSetSync(b *testing.B) {
	dir := b.TempDir()
	s := New(dir, config.DefaultConfig().WithSync())
	if err := s.Init(); err != nil {
		b.Error(err)
	}
	defer s.Close()

	tests := []benchmarkTestCase{
		{"128B", 128},
		{"256B", 256},
		{"1K", 1024},
		{"2K", 2048},
		{"4K", 4096},
		{"8K", 8192},
		{"16K", 16384},
		{"32K", 32768},
		{"64K", 65536},
		{"128K", 131072},
		{"256K", 262144},
		{"512K", 524288},
		{"1M", 1048576},
		{"2M", 2097152},
		{"4M", 4194304},
		{"8M", 8388608},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.SetBytes(int64(test.size))

			key := "foo"
			value := []byte(strings.Repeat(" ", test.size))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := s.Set([]byte(key), value); err != nil {
					b.Error(err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkStupidSetAsyncSync(b *testing.B) {
	dir := b.TempDir()
	s := New(dir, config.DefaultConfig().WithAsyncSync())
	if err := s.Init(); err != nil {
		b.Error(err)
	}
	defer s.Close()

	tests := []benchmarkTestCase{
		{"128B", 128},
		{"256B", 256},
		{"1K", 1024},
		{"2K", 2048},
		{"4K", 4096},
		{"8K", 8192},
		{"16K", 16384},
		{"32K", 32768},
		{"64K", 65536},
		{"128K", 131072},
		{"256K", 262144},
		{"512K", 524288},
		{"1M", 1048576},
		{"2M", 2097152},
		{"4M", 4194304},
		{"8M", 8388608},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.SetBytes(int64(test.size))

			key := "foo"
			value := []byte(strings.Repeat(" ", test.size))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if err := s.Set([]byte(key), value); err != nil {
					b.Error(err)
				}
			}
			b.StopTimer()
		})
	}
}

func BenchmarkStupidGet(b *testing.B) {
	tests := []benchmarkTestCase{
		{"128B", 128},
		{"256B", 256},
		{"1K", 1024},
		{"2K", 2048},
		{"4K", 4096},
		{"8K", 8192},
		{"16K", 16384},
		{"32K", 32768},
		{"64K", 65536},
		{"128K", 131072},
		{"256K", 262144},
		{"512K", 524288},
		{"1M", 1048576},
		{"2M", 2097152},
		{"4M", 4194304},
		{"8M", 8388608},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			dir := b.TempDir()
			s := New(dir, config.DefaultConfig().WithSync())
			if err := s.Init(); err != nil {
				b.Error(err)
			}
			defer s.Close()

			b.SetBytes(int64(test.size))

			key := "foo"
			value := []byte(strings.Repeat(" ", test.size))

			if err := s.Set([]byte(key), value); err != nil {
				b.Error(err)
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if val, err := s.Get([]byte(key)); err != nil {
					b.Error(err)
				} else if !bytes.Equal(val, value) {
					b.Error("expected", value, "got", val)
				}
			}
			b.StopTimer()
		})
	}
}
