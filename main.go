package main

import (
	"os"
	"path/filepath"

	"github.com/utkarsh-pro/use/pkg/storage"
)

func main() {
	dbpath := filepath.Join(".", "db")
	if err := os.MkdirAll(dbpath, 0777); err != nil {
		panic(err)
	}

	s, err := storage.New(storage.StupidStorageType, dbpath)
	if err != nil {
		panic(err)
	}

	if err := s.Init(); err != nil {
		panic(err)
	}

	if err := s.Set("foo", []byte("bar")); err != nil {
		panic(err)
	}

	if err := s.Set("foo", []byte("bazz")); err != nil {
		panic(err)
	}

	val, err := s.Get("foo")
	if err != nil {
		panic(err)
	}
	println(string(val))

	// if err := s.Delete("foo"); err != nil {
	// 	panic(err)
	// }

	ok, err := s.Exists("foo")
	if err != nil {
		panic(err)
	}
	println(ok)
}
