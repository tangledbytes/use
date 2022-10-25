package main

import (
	"github.com/utkarsh-pro/use/pkg/config"
	"github.com/utkarsh-pro/use/pkg/log"
	"github.com/utkarsh-pro/use/pkg/shutdown"
	"github.com/utkarsh-pro/use/pkg/storage"
	scfg "github.com/utkarsh-pro/use/pkg/storage/config"
	"github.com/utkarsh-pro/use/pkg/transport"
)

func main() {
	config.Setup()
	log.SetLevel(config.LogLevel)

	storage, err := storage.New(storage.StorageType(config.Storage), config.StoragePath, scfg.DefaultConfig())
	if err != nil {
		panic(err)
	}

	if err := storage.Init(); err != nil {
		panic(err)
	}

	transport, err := transport.New(transport.TransportType(config.Transport), storage)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := transport.Setup(config.Address); err != nil {
			panic(err)
		}
	}()

	log.Infoln("use is running: ", config.String())

	shutdown.RegisterFunc(storage.Close)
	shutdown.RegisterFunc(transport.Shutdown)

	shutdown.OnSignal()
}
